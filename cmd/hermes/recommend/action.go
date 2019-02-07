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

type Forecast struct {
	AccountID      string          `json:"account_id"`
	Alias          string          `json:"alias"`
	Region         string          `json:"region"`
	UsageType      string          `json:"usage_type"`
	Platform       string          `json:"platform,omitempty"`
	CacheEngine    string          `json:"cache_engine,omitempty"`
	DatabaseEngine string          `json:"database_engine,omitempty"`
	InstanceNum    InstanceNumList `json:"instance_num"`
	Extend         interface{}     `json:"extend,omitempty"`
}

type ForecastList struct {
	Forecast []*Forecast `json:"forecast"`
}

func (list ForecastList) JSON() string {
	bytea, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

func (list ForecastList) Merge() MergedForecastList {
	flat := make(map[string]InstanceNumList)
	for _, in := range list.Forecast {
		key := fmt.Sprintf("%s_%s_%s_%s_%s",
			in.Region,
			in.UsageType,
			in.Platform,
			in.CacheEngine,
			in.DatabaseEngine,
		)

		val, ok := flat[key]
		if ok {
			flat[key] = val.Add(in.InstanceNum)
			continue
		}

		flat[key] = in.InstanceNum
	}

	out := MergedForecastList{}
	for k, v := range flat {
		keys := strings.Split(k, "_")
		out = append(out, &MergedForecast{
			Region:         keys[0],
			UsageType:      keys[1],
			Platform:       keys[2],
			CacheEngine:    keys[3],
			DatabaseEngine: keys[4],
			InstanceNum:    v,
		})
	}

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].UsageType < out[j].UsageType
	})
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Platform < out[j].Platform
	})

	return out
}

func (list ForecastList) Load() {
	flat := make(map[string]bool)
	for _, f := range list.Forecast {
		flat[f.Region] = true
	}

	region := []string{}
	for k := range flat {
		region = append(region, k)
	}

	path := fmt.Sprintf("%s/pricing", tmpdir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(tmpdir, os.ModePerm)
	}

	for _, r := range region {
		cache := fmt.Sprintf("%s/%s.out", path, r)
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

type InstanceNum struct {
	Date        string  `json:"date"`
	InstanceNum float64 `json:"instance_num"`
}

type InstanceNumList []InstanceNum

func (list InstanceNumList) Add(input InstanceNumList) InstanceNumList {
	out := InstanceNumList{}
	for i := range list {
		if list[i].Date != input[i].Date {
			panic(fmt.Sprintf("invalid args %v %v", list[i], input[i]))
		}

		out = append(out, InstanceNum{
			Date:        list[i].Date,
			InstanceNum: list[i].InstanceNum + input[i].InstanceNum,
		})
	}

	return out
}

type MergedForecast struct {
	Region         string        `json:"region"`
	UsageType      string        `json:"usage_type"`
	Platform       string        `json:"platform,omitempty"`
	CacheEngine    string        `json:"cache_engine,omitempty"`
	DatabaseEngine string        `json:"database_engine,omitempty"`
	InstanceNum    []InstanceNum `json:"instance_num"`
}

type MergedForecastList []*MergedForecast

func (list MergedForecastList) Recommended() (pricing.RecommendedList, error) {
	out := pricing.RecommendedList{}
	for _, in := range list {
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

		os := pricing.OperatingSystem[in.Platform]
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

	for _, in := range list {
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

	for _, in := range list {
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

type Output struct {
	Forecast       ForecastList            `json:"forecast"`
	MergedForecast MergedForecastList      `json:"merged_forecast"`
	Recommended    pricing.RecommendedList `json:"recommended"`
	Result         pricing.RecommendedList `json:"result"`
}

func (output *Output) Array() [][]interface{} {
	array := [][]interface{}{}

	forecast := []interface{}{
		"forecast", "account_id", "alies", "usage_type", "platform/engine",
	}
	for _, n := range output.Forecast.Forecast[0].InstanceNum {
		forecast = append(forecast, n.Date)
	}
	array = append(array, forecast)

	for _, f := range output.Forecast.Forecast {
		val := []interface{}{""}

		val = append(val, f.AccountID)
		val = append(val, f.Alias)
		val = append(val, f.UsageType)
		if len(f.Platform) > 0 {
			val = append(val, f.Platform)
		}
		if len(f.DatabaseEngine) > 0 {
			val = append(val, f.DatabaseEngine)
		}
		if len(f.CacheEngine) > 0 {
			val = append(val, f.CacheEngine)
		}
		for _, n := range f.InstanceNum {
			val = append(val, n.InstanceNum)
		}

		array = append(array, val)
	}
	array = append(array, []interface{}{""})

	merged := []interface{}{
		"merged_forecast", "", "", "usage_type", "platform/engine",
	}
	for _, n := range output.Forecast.Forecast[0].InstanceNum {
		merged = append(merged, n.Date)
	}
	array = append(array, merged)

	for _, m := range output.MergedForecast {
		val := []interface{}{"", "", ""}

		val = append(val, m.UsageType)
		if len(m.Platform) > 0 {
			val = append(val, m.Platform)
		}
		if len(m.DatabaseEngine) > 0 {
			val = append(val, m.DatabaseEngine)
		}
		if len(m.CacheEngine) > 0 {
			val = append(val, m.CacheEngine)
		}
		for _, n := range m.InstanceNum {
			val = append(val, n.InstanceNum)
		}

		array = append(array, val)
	}
	array = append(array, []interface{}{""})

	recommended := []interface{}{
		"recommended", "", "", "usage_type", "os/engine", "ondemand_num_avg", "reserved_num", "full_ondemand_cost", "reserved_applied_cost.ondemand", "reserved_applied_cost.reserved", "reserved_applied_cost.total", "subtraction", "discount_rate", "reserved_quantity",
	}
	array = append(array, recommended)

	for _, r := range output.Recommended {
		val := []interface{}{"", "", ""}

		val = append(val, r.Record.UsageType)
		if len(r.Record.OperatingSystem) > 0 {
			val = append(val, r.Record.OperatingSystem)
		}
		if len(r.Record.CacheEngine) > 0 {
			val = append(val, r.Record.CacheEngine)
		}
		if len(r.Record.DatabaseEngine) > 0 {
			val = append(val, r.Record.DatabaseEngine)
		}

		val = append(val, r.OnDemandInstanceNumAvg)
		val = append(val, r.ReservedInstanceNum)
		val = append(val, r.FullOnDemandCost)
		val = append(val, r.ReservedAppliedCost.OnDemand)
		val = append(val, r.ReservedAppliedCost.Reserved)
		val = append(val, r.ReservedAppliedCost.Total)
		val = append(val, r.Subtraction)
		val = append(val, r.DiscountRate)
		val = append(val, r.ReservedQuantity)

		array = append(array, val)
	}
	array = append(array, []interface{}{""})

	result := []interface{}{
		"result", "", "", "usage_type", "os/engine", "instance_num",
	}
	array = append(array, result)

	for _, r := range output.Result {
		val := []interface{}{"", "", ""}

		val = append(val, r.MinimumRecord.UsageType)
		if len(r.MinimumRecord.OperatingSystem) > 0 {
			val = append(val, r.MinimumRecord.OperatingSystem)
		}
		if len(r.MinimumRecord.DatabaseEngine) > 0 {
			val = append(val, r.MinimumRecord.DatabaseEngine)
		}
		if len(r.MinimumRecord.CacheEngine) > 0 {
			val = append(val, r.MinimumRecord.CacheEngine)
		}

		val = append(val, r.MinimumReservedInstanceNum)

		array = append(array, val)
	}

	return array
}

func Action(c *cli.Context) {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(fmt.Errorf("stdin: %v", err))
		return
	}

	var input ForecastList
	if uerr := json.Unmarshal(stdin, &input); uerr != nil {
		fmt.Println(fmt.Errorf("unmarshal: %v", uerr))
		return
	}

	input.Load()
	mf := input.Merge()
	recommended, err := mf.Recommended()
	if err != nil {
		fmt.Println(fmt.Errorf("recommended: %v", err))
		return
	}

	output := &Output{
		Forecast:       input,
		MergedForecast: mf,
		Recommended:    recommended,
		Result:         recommended.Merge(),
	}

	if c.String("format") == "csv" {
		for _, r := range output.Array() {
			for _, c := range r {
				fmt.Printf("%v, ", c)
			}
			fmt.Println()
		}
		return
	}

	if c.String("format") == "tsv" {
		for _, r := range output.Array() {
			for _, c := range r {
				fmt.Printf("%v\t", c)
			}
			fmt.Println()
		}
		return
	}

	//  c.String("format") == "json"
	bytes, err := json.Marshal(&output)
	if err != nil {
		fmt.Println(fmt.Errorf("marshal: %v", err))
		return
	}

	fmt.Println(string(bytes))
}
