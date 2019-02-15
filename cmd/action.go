package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/reservation"
	"github.com/urfave/cli"
)

var tmpdir = "/var/tmp/hermes"

type Input struct {
	Forecast ForecastList `json:"forecast"`
}

func (input *Input) JSON() string {
	bytea, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

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

func (f *Forecast) PlatformEngine() string {
	if len(f.Platform) > 0 {
		return f.Platform
	}

	if len(f.DatabaseEngine) > 0 {
		return f.DatabaseEngine
	}

	if len(f.CacheEngine) > 0 {
		return f.CacheEngine
	}

	return ""
}

type ForecastList []*Forecast

func (list ForecastList) Merge() MergedForecastList {
	flat := make(map[string]InstanceNumList)
	for _, in := range list {
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

	return out
}

func (list ForecastList) Load() {
	flat := make(map[string]bool)
	for _, f := range list {
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

	{
		path := fmt.Sprintf("%s/reservation.out", tmpdir)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			repo := reservation.NewRepository(region)
			if err := repo.Fetch(); err != nil {
				fmt.Println(fmt.Errorf("fetch reservation: %v", err))
				return
			}

			if err := repo.Write(path); err != nil {
				fmt.Println(fmt.Errorf("write reservation: %v", err))
				return
			}
		}
	}
}

func (list ForecastList) Array() [][]interface{} {
	array := [][]interface{}{}

	header := []interface{}{
		"account_id",
		"alias",
		"usage_type",
		"platform/engine",
	}
	for _, n := range list[0].InstanceNum {
		header = append(header, n.Date)
	}
	array = append(array, header)

	for _, f := range list {
		val := []interface{}{
			f.AccountID,
			f.Alias,
			f.UsageType,
			f.PlatformEngine(),
		}

		for _, n := range f.InstanceNum {
			val = append(val, n.InstanceNum)
		}

		array = append(array, val)
	}

	return array
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
	Region         string          `json:"region"`
	UsageType      string          `json:"usage_type"`
	Platform       string          `json:"platform,omitempty"`
	CacheEngine    string          `json:"cache_engine,omitempty"`
	DatabaseEngine string          `json:"database_engine,omitempty"`
	InstanceNum    InstanceNumList `json:"instance_num"`
}

func (f *MergedForecast) PlatformEngine() string {
	if len(f.Platform) > 0 {
		return f.Platform
	}

	if len(f.DatabaseEngine) > 0 {
		return f.DatabaseEngine
	}

	if len(f.CacheEngine) > 0 {
		return f.CacheEngine
	}

	panic("platform/engine notfound")
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

func (list MergedForecastList) Array(date []string) [][]interface{} {
	array := [][]interface{}{}

	header := []interface{}{
		"account_id",
		"alias",
		"usage_type",
		"platform/engine",
	}
	for _, d := range date {
		header = append(header, d)
	}
	array = append(array, header)

	for _, m := range list {
		val := []interface{}{
			"n/a",
			"n/a",
			m.UsageType,
			m.PlatformEngine(),
		}

		for _, n := range m.InstanceNum {
			val = append(val, n.InstanceNum)
		}

		array = append(array, val)
	}

	return array
}

type Result struct {
	UsageType   string  `json:"usage_type"`
	OSEngine    string  `json:"os_engine"`
	InstanceNum float64 `json:"instance_num"`
	CurrentRI   float64 `json:"current_ri"`
	Difference  float64 `json:"difference"`
}

type ResultList []*Result

func (list ResultList) Array() [][]interface{} {
	array := append([][]interface{}{}, []interface{}{
		"usage_type",
		"os/engine",
		"instance_num",
		"current_ri",
		"difference",
	})

	for _, r := range list {
		array = append(array, []interface{}{
			r.UsageType,
			r.OSEngine,
			r.InstanceNum,
			r.CurrentRI,
			r.Difference,
		})
	}

	return array
}

func NewResultList(list pricing.RecommendedList) (ResultList, error) {
	out := ResultList{}

	repo, err := reservation.Read("/var/tmp/hermes/reservation.out")
	if err != nil {
		return nil, fmt.Errorf("read reservation: %v", err)
	}

	m := list.Merge()
	for _, r := range m {
		min := r.MinimumRecord
		if len(min.OperatingSystem) < 1 {
			continue
		}

		rs := repo.FindByInstanceType(min.InstanceType).
			Region(min.Region).
			LeaseContractLength(min.LeaseContractLength).
			OfferingClass(min.OfferingClass).
			OfferingType(min.PurchaseOption).
			ProductDescription(min.OSEngine()).
			Active()

		var current float64
		if len(rs) == 1 {
			current = float64(rs[0].Count())
		} else if len(rs) == 0 {
			// not found.
		} else {
			return nil, fmt.Errorf("invalid ec2 reservation: %v", rs)
		}

		out = append(out, &Result{
			UsageType:   min.UsageType,
			OSEngine:    min.OSEngine(),
			InstanceNum: r.MinimumReservedInstanceNum,
			CurrentRI:   current,
			Difference:  r.MinimumReservedInstanceNum - current,
		})
	}

	for _, r := range m {
		min := r.MinimumRecord
		if len(min.CacheEngine) < 1 {
			continue
		}

		rs := repo.SelectAll().
			CacheNodeType(min.InstanceType).
			Region(min.Region).
			LeaseContractLength(min.LeaseContractLength).
			OfferingType(min.PurchaseOption).
			ProductDescription(min.OSEngine()).
			Active()

		var current float64
		if len(rs) == 1 {
			current = float64(rs[0].Count())
		} else if len(rs) == 0 {
			// not found.
		} else {
			return nil, fmt.Errorf("invalid cache reservation: %v", rs)
		}

		out = append(out, &Result{
			UsageType:   min.UsageType,
			OSEngine:    min.OSEngine(),
			InstanceNum: r.MinimumReservedInstanceNum,
			CurrentRI:   current,
			Difference:  r.MinimumReservedInstanceNum - current,
		})
	}

	for _, r := range m {
		min := r.MinimumRecord
		if len(min.DatabaseEngine) < 1 {
			continue
		}

		maz := false
		if strings.Contains(min.UsageType, "Multi-AZ") {
			maz = true
		}

		rs := repo.SelectAll().
			DBInstanceClass(min.InstanceType).
			Region(min.Region).
			LeaseContractLength(min.LeaseContractLength).
			OfferingType(min.PurchaseOption).
			ProductDescription(min.OSEngine()).
			MultiAZ(maz).
			Active()

		var current float64
		if len(rs) == 1 {
			current = float64(rs[0].Count())
		} else if len(rs) == 0 {
			// not found.
		} else {
			return nil, fmt.Errorf("invalid database reservation: %v", rs)
		}

		out = append(out, &Result{
			UsageType:   min.UsageType,
			OSEngine:    min.OSEngine(),
			InstanceNum: r.MinimumReservedInstanceNum,
			CurrentRI:   current,
			Difference:  r.MinimumReservedInstanceNum - current,
		})
	}

	return out, nil
}

type Output struct {
	Forecast       ForecastList            `json:"forecast"`
	MergedForecast MergedForecastList      `json:"merged_forecast"`
	Recommended    pricing.RecommendedList `json:"recommended"`
	Result         ResultList              `json:"result"`
}

func (output *Output) Array() [][]interface{} {
	array := [][]interface{}{}

	array = append(array, output.Forecast.Array()...)
	array = append(array, []interface{}{""})

	date := []string{}
	for _, d := range output.Forecast[0].InstanceNum {
		date = append(date, d.Date)
	}

	array = append(array, output.MergedForecast.Array(date)...)
	array = append(array, []interface{}{""})

	array = append(array, output.Recommended.Array()...)
	array = append(array, []interface{}{""})

	array = append(array, output.Result.Array()...)

	return array
}

func Action(c *cli.Context) {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(fmt.Errorf("stdin: %v", err))
		return
	}

	var input Input
	if uerr := json.Unmarshal(stdin, &input); uerr != nil {
		fmt.Println(fmt.Errorf("unmarshal: %v", uerr))
		return
	}

	input.Forecast.Load()

	mf := input.Forecast.Merge()
	rec, err := mf.Recommended()
	if err != nil {
		fmt.Println(fmt.Errorf("recommended: %v", err))
		return
	}
	res, err := NewResultList(rec)
	if err != nil {
		fmt.Println(fmt.Errorf("new result list: %v", err))
		return
	}

	output := &Output{
		Forecast:       input.Forecast,
		MergedForecast: mf,
		Recommended:    rec,
		Result:         res,
	}

	if c.String("format") == "csv" {
		for _, r := range output.Array() {
			for _, c := range r {
				fmt.Printf("%v,", c)
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
