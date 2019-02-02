package recommend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/urfave/cli"
)

var tmpdir = "/var/tmp/hermes"

type ForecstList struct {
	Forecast []*Forecast `json:"forecast"`
}

type Forecast struct {
	AccountID      string        `json:"account_id"`
	Alias          string        `json:"alias"`
	Region         string        `json:"region"`
	UsageType      string        `json:"usage_type"`
	Platform       string        `json:"platform,omitempty"`
	CacheEngine    string        `json:"cache_engine,omitempty"`
	DatabaseEngine string        `json:"database_engine,omitempty"`
	InstanceNum    []InstanceNum `json:"instance_num"`
	Extend         interface{}   `json:"extend,omitempty"`
}

type InstanceNum struct {
	Date        string  `json:"date"`
	InstanceNum float64 `json:"instance_num"`
}

type Output struct {
	Forecast    []*Forecast            `json:"forecast"`
	Merged      []*Merged              `json:"merged"`
	Recommended []*pricing.Recommended `json:"recommended"`
}

type Merged struct {
	Region         string        `json:"region"`
	UsageType      string        `json:"usage_type"`
	Platform       string        `json:"platform,omitempty"`
	CacheEngine    string        `json:"cache_engine,omitempty"`
	DatabaseEngine string        `json:"database_engine,omitempty"`
	InstanceNum    []InstanceNum `json:"instance_num"`
}

func (input ForecstList) JSON() string {
	bytea, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

func Action(c *cli.Context) {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(fmt.Errorf("stdin: %v", err))
		return
	}

	var input ForecstList
	if uerr := json.Unmarshal(stdin, &input); uerr != nil {
		fmt.Println(fmt.Errorf("unmarshal: %v", uerr))
		return
	}

	Cache(input.Forecast)
	merged := Merge(input.Forecast)
	recommended, err := Recommended(merged)
	if err != nil {
		fmt.Println(fmt.Errorf("recommended: %v", err))
		return
	}

	output := Output{
		Forecast:    input.Forecast,
		Merged:      merged,
		Recommended: recommended,
	}

	bytes, err := json.Marshal(&output)
	if err != nil {
		fmt.Println(fmt.Errorf("marshal: %v", err))
		return
	}

	fmt.Println(string(bytes))
}

func Recommended(merged []*Merged) ([]*pricing.Recommended, error) {
	out := []*pricing.Recommended{}
	for _, in := range merged {
		if len(in.Platform) < 1 {
			continue
		}

		forecast := []pricing.Forecast{}
		for _, n := range in.InstanceNum {
			forecast = append(forecast, pricing.Forecast{
				Date:        n.Date,
				InstanceNum: n.InstanceNum,
			})
		}

		os := pricing.OperationgSystem[in.Platform]
		path := fmt.Sprintf("%s/pricing/%s.out", tmpdir, in.Region)
		repo, err := pricing.Read(path)
		if err != nil {
			return nil, fmt.Errorf("read pricing (region=%s): %v", in.Region, err)
		}

		hit := repo.FindByUsageType(in.UsageType).
			OperatingSystem(os).
			LeaseContractLength("1yr").
			PurchaseOption("All Upfront").
			OfferingClass("standard").
			PreInstalled("NA").
			Tenancy("Shared")

		if len(hit) != 1 {
			continue
		}

		recommend, err := repo.Recommend(hit[0], forecast)
		if err != nil {
			return nil, fmt.Errorf("recommend ec2: %v", err)
		}

		if recommend.ReservedInstanceNum > 0 {
			out = append(out, recommend)
		}
	}

	for _, in := range merged {
		if len(in.CacheEngine) < 1 {
			continue
		}

		forecast := []pricing.Forecast{}
		for _, n := range in.InstanceNum {
			forecast = append(forecast, pricing.Forecast{
				Date:        n.Date,
				InstanceNum: n.InstanceNum,
			})
		}

		path := fmt.Sprintf("%s/pricing/%s.out", tmpdir, in.Region)
		repo, err := pricing.Read(path)
		if err != nil {
			return nil, fmt.Errorf("read pricing (region=%s): %v", in.Region, err)
		}

		hit := repo.FindByUsageType(in.UsageType).
			LeaseContractLength("1yr").
			PurchaseOption("Heavy Utilization").
			CacheEngine(in.CacheEngine)

		if len(hit) != 1 {
			continue
		}

		recommend, err := repo.Recommend(hit[0], forecast, "minimum")
		if err != nil {
			return nil, fmt.Errorf("recommend cache: %v", err)
		}

		if recommend.ReservedInstanceNum > 0 {
			out = append(out, recommend)
		}
	}

	for _, in := range merged {
		if len(in.DatabaseEngine) < 1 {
			continue
		}

		forecast := []pricing.Forecast{}
		for _, n := range in.InstanceNum {
			forecast = append(forecast, pricing.Forecast{
				Date:        n.Date,
				InstanceNum: n.InstanceNum,
			})
		}

		path := fmt.Sprintf("%s/pricing/%s.out", tmpdir, in.Region)
		repo, err := pricing.Read(path)
		if err != nil {
			return nil, fmt.Errorf("read pricing (region=%s): %v", in.Region, err)
		}

		hit := repo.FindByUsageType(in.UsageType).
			LeaseContractLength("1yr").
			PurchaseOption("All Upfront").
			DatabaseEngine(in.DatabaseEngine)

		if len(hit) != 1 {
			continue
		}

		recommend, err := repo.Recommend(hit[0], forecast, "minimum")
		if err != nil {
			return nil, fmt.Errorf("recommend rds: %v", err)
		}

		if recommend.ReservedInstanceNum > 0 {
			out = append(out, recommend)
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Record.UsageType < out[j].Record.UsageType
	})

	return out, nil
}

func Merge(forecast []*Forecast) []*Merged {
	flat := make(map[string][]InstanceNum)
	for _, in := range forecast {
		key := fmt.Sprintf("%s_%s_%s_%s_%s",
			in.Region,
			in.UsageType,
			in.Platform,
			in.CacheEngine,
			in.DatabaseEngine,
		)

		val, ok := flat[key]
		if ok {
			flat[key] = Add(val, in.InstanceNum)
			continue
		}

		flat[key] = in.InstanceNum
	}

	out := []*Merged{}
	for k, v := range flat {
		keys := strings.Split(k, "_")
		out = append(out, &Merged{
			Region:         keys[0],
			UsageType:      keys[1],
			Platform:       keys[2],
			CacheEngine:    keys[3],
			DatabaseEngine: keys[4],
			InstanceNum:    v,
		})
	}

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Region < out[j].Region
	})
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].UsageType < out[j].UsageType
	})
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Platform < out[j].Platform
	})

	return out
}

func Add(val []InstanceNum, input []InstanceNum) []InstanceNum {
	list := []InstanceNum{}
	for i := range val {
		if val[i].Date != input[i].Date {
			panic(fmt.Sprintf("invalid args %v %v", val[i], input[i]))
		}

		list = append(list, InstanceNum{
			Date:        val[i].Date,
			InstanceNum: val[i].InstanceNum + input[i].InstanceNum,
		})
	}

	return list
}

func Cache(forecast []*Forecast) {
	rflat := make(map[string]bool)
	for _, f := range forecast {
		rflat[f.Region] = true
	}

	region := []string{}
	for k, _ := range rflat {
		region = append(region, k)
	}

	for _, r := range region {
		cache := fmt.Sprintf("%s/pricing/%s.out", tmpdir, r)
		if _, err := os.Stat(cache); os.IsNotExist(err) {
			repo := pricing.NewRepository([]string{r})
			if err := repo.Fetch(); err != nil {
				fmt.Println(fmt.Errorf("fetch pricing (region=%s): %v", r, err))
				return
			}

			if err := repo.Write(cache); err != nil {
				fmt.Println(fmt.Errorf("write pricing (region=%s): %v", r, err))
				return
			}
		}
	}
}
