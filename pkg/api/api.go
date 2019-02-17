package api

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/reserved"
)

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

type Output struct {
	Forecast    ForecastList            `json:"forecast"`
	Merged      ForecastList            `json:"merged"`
	Recommended pricing.RecommendedList `json:"recommended"`
	Coverage    CoverageList            `json:"coverage"`
}

func (output *Output) Array() [][]interface{} {
	array := [][]interface{}{}

	array = append(array, []interface{}{"# forecast usage"})
	array = append(array, output.Forecast.Array()...)

	array = append(array, []interface{}{"# forecast usage merged"})
	array = append(array, output.Merged.Array()...)

	array = append(array, []interface{}{"# recommended reserved instances"})
	array = append(array, output.Recommended.Array()...)

	array = append(array, []interface{}{"# coverage"})
	array = append(array, output.Coverage.Array()...)

	return array
}

type Forecast struct {
	AccountID      string          `json:"account_id,omitempty"`
	Alias          string          `json:"alias,omitempty"`
	Region         string          `json:"region"`
	UsageType      string          `json:"usage_type"`
	Platform       string          `json:"platform,omitempty"`
	CacheEngine    string          `json:"cache_engine,omitempty"`
	DatabaseEngine string          `json:"database_engine,omitempty"`
	InstanceNum    InstanceNumList `json:"instance_num"`
	Extend         interface{}     `json:"extend,omitempty"`
}

func (f *Forecast) Date() []string {
	date := []string{}
	for _, d := range f.InstanceNum {
		date = append(date, d.Date)
	}
	return date
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

	panic("platform/engine not found")
}

type ForecastList []*Forecast

func (list ForecastList) Merge() ForecastList {
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

	out := ForecastList{}
	for k, v := range flat {
		keys := strings.Split(k, "_")
		out = append(out, &Forecast{
			AccountID:      "n/a",
			Alias:          "n/a",
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

type InstanceNumList []*InstanceNum

func (list InstanceNumList) Add(input InstanceNumList) InstanceNumList {
	out := InstanceNumList{}
	for i := range list {
		if list[i].Date != input[i].Date {
			panic(fmt.Sprintf("invalid args %v %v", list[i], input[i]))
		}

		out = append(out, &InstanceNum{
			Date:        list[i].Date,
			InstanceNum: list[i].InstanceNum + input[i].InstanceNum,
		})
	}

	return out
}

func Recommend(list ForecastList, repo []*pricing.Repository) (pricing.RecommendedList, error) {
	rmap := make(map[string]*pricing.Repository)
	for i := range repo {
		rmap[repo[i].Region[0]] = repo[i]
	}

	out := pricing.RecommendedList{}
	for _, in := range list {
		if len(in.Platform) < 1 {
			continue
		}

		forecast := pricing.ForecastList{}
		for _, n := range in.InstanceNum {
			forecast = append(forecast, &pricing.Forecast{
				Date:        n.Date,
				InstanceNum: n.InstanceNum,
			})
		}

		repo := rmap[in.Region]
		hit := repo.FindByUsageType(in.UsageType).
			OperatingSystem(pricing.OperatingSystem[in.Platform]).
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

		forecast := pricing.ForecastList{}
		for _, n := range in.InstanceNum {
			forecast = append(forecast, &pricing.Forecast{
				Date:        n.Date,
				InstanceNum: n.InstanceNum,
			})
		}

		repo := rmap[in.Region]
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

		forecast := pricing.ForecastList{}
		for _, n := range in.InstanceNum {
			forecast = append(forecast, &pricing.Forecast{
				Date:        n.Date,
				InstanceNum: n.InstanceNum,
			})
		}

		repo := rmap[in.Region]
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

type Coverage struct {
	UsageType   string  `json:"usage_type"`
	OSEngine    string  `json:"os_engine"`
	InstanceNum float64 `json:"instance_num"`
	CurrentRI   float64 `json:"current_ri"`
	Short       float64 `json:"short"`
	Coverage    float64 `json:"coverage"`
}

type CoverageList []*Coverage

func (list CoverageList) Array() [][]interface{} {
	array := append([][]interface{}{}, []interface{}{
		"usage_type",
		"os/engine",
		"instance_num",
		"current_ri",
		"short",
		"coverage",
	})

	for _, r := range list {
		array = append(array, []interface{}{
			r.UsageType,
			r.OSEngine,
			r.InstanceNum,
			r.CurrentRI,
			r.Short,
			r.Coverage,
		})
	}

	return array
}

func GetCoverage(list pricing.RecommendedList, repo *reserved.Repository) (CoverageList, error) {
	used := reserved.RecordList{}
	out := CoverageList{}
	for i := range list {
		if !list[i].MinimumRecord.IsInstance() {
			continue
		}

		min := list[i].MinimumRecord
		rs := repo.FindByInstanceType(min.InstanceType).
			Region(min.Region).
			LeaseContractLength(min.LeaseContractLength).
			OfferingClass(min.OfferingClass).
			OfferingType(min.PurchaseOption).
			ProductDescription(min.OSEngine()).
			Active()

		var current float64
		if len(rs) == 0 {
			// not found
		} else if len(rs) > 0 {
			current = float64(rs.CountSum())
			used = append(used, rs...)
		} else {
			return nil, fmt.Errorf("invalid ec2 reservation: %v", rs)
		}

		out = append(out, &Coverage{
			UsageType:   min.UsageType,
			OSEngine:    min.OSEngine(),
			InstanceNum: list[i].MinimumReservedInstanceNum,
			CurrentRI:   current,
			Short:       list[i].MinimumReservedInstanceNum - current,
			Coverage:    current / list[i].MinimumReservedInstanceNum,
		})
	}

	for i := range list {
		if !list[i].MinimumRecord.IsCacheNode() {
			continue
		}

		min := list[i].MinimumRecord
		rs := repo.SelectAll().
			CacheNodeType(min.InstanceType).
			Region(min.Region).
			LeaseContractLength(min.LeaseContractLength).
			OfferingType(min.PurchaseOption).
			ProductDescription(min.OSEngine()).
			Active()

		var current float64
		if len(rs) == 0 {
			// not found
		} else if len(rs) > 0 {
			current = float64(rs.CountSum())
			used = append(used, rs...)
		} else {
			return nil, fmt.Errorf("invalid cache reservation: %v", rs)
		}

		out = append(out, &Coverage{
			UsageType:   min.UsageType,
			OSEngine:    min.OSEngine(),
			InstanceNum: list[i].MinimumReservedInstanceNum,
			CurrentRI:   current,
			Short:       list[i].MinimumReservedInstanceNum - current,
			Coverage:    current / list[i].MinimumReservedInstanceNum,
		})
	}

	for i := range list {
		if !list[i].MinimumRecord.IsDatabase() {
			continue
		}

		min := list[i].MinimumRecord
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
		if len(rs) == 0 {
			// not found
		} else if len(rs) > 0 {
			current = float64(rs.CountSum())
			used = append(used, rs...)
		} else {
			return nil, fmt.Errorf("invalid database reservation: %v", rs)
		}

		out = append(out, &Coverage{
			UsageType:   min.UsageType,
			OSEngine:    min.OSEngine(),
			InstanceNum: list[i].MinimumReservedInstanceNum,
			CurrentRI:   current,
			Short:       list[i].MinimumReservedInstanceNum - current,
			Coverage:    current / list[i].MinimumReservedInstanceNum,
		})
	}

	unused := reserved.RecordList{}
	for _, r := range repo.SelectAll().Active() {
		found := false
		for _, u := range used {
			if r.Equals(u) {
				found = true
			}
		}

		if !found {
			unused = append(unused, r)
		}
	}

	for _, r := range unused {
		out = append(out, &Coverage{
			UsageType:   r.Region + "-" + r.InstanceType + r.CacheNodeType + r.DBInstanceClass,
			OSEngine:    r.ProductDescription,
			InstanceNum: 0,
			CurrentRI:   float64(r.Count()),
			Short:       float64(-r.Count()),
			Coverage:    math.MaxFloat64,
		})
	}

	return out, nil
}
