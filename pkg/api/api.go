package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/itsubaki/hermes/internal/costexp"
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

func (output *Output) CSV() string {
	var buf bytes.Buffer
	for _, r := range output.Array() {
		for _, c := range r {
			buf.WriteString(fmt.Sprintf("%v,", c))
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

func (output *Output) TSV() string {
	var buf bytes.Buffer
	for _, r := range output.Array() {
		for _, c := range r {
			buf.WriteString(fmt.Sprintf("%v\t", c))
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

func (output *Output) JSON() string {
	bytea, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

func (output *Output) Array() [][]interface{} {
	array := [][]interface{}{}

	array = append(array, []interface{}{"# forecast usage"})
	array = append(array, output.Forecast.Header())
	array = append(array, output.Forecast.Array()...)

	array = append(array, []interface{}{"# forecast usage merged"})
	array = append(array, output.Merged.Header())
	array = append(array, output.Merged.Array()...)

	array = append(array, []interface{}{"# recommended reserved instances"})
	array = append(array, output.Recommended.Header())
	array = append(array, output.Recommended.Array()...)

	array = append(array, []interface{}{"# coverage"})
	array = append(array, output.Coverage.Header())
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

func (list ForecastList) RecommendCompute(repo []*pricing.Repository) (pricing.RecommendedList, error) {
	rmap := make(map[string]*pricing.Repository)
	for i := range repo {
		rmap[repo[i].Region[0]] = repo[i]
	}

	out := pricing.RecommendedList{}
	for _, in := range list {
		if len(in.Platform) < 1 {
			continue
		}

		repo := rmap[in.Region]
		price := repo.SelectAll().
			UsageType(in.UsageType).
			OperatingSystem(pricing.OperatingSystem[in.Platform]).
			LeaseContractLength("1yr").
			PurchaseOption("All Upfront").
			OfferingClass("standard").
			PreInstalled("NA").
			Tenancy("Shared")

		if len(price) != 1 {
			continue
		}

		forecast := in.InstanceNum.ForecastList()
		rec, err := repo.Recommend(price[0], forecast)
		if err != nil {
			return nil, fmt.Errorf("recommend compute: %v", err)
		}

		if rec.ReservedInstanceNum > 0 {
			out = append(out, rec)
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Record.UsageType < out[j].Record.UsageType
	})

	return out, nil
}

func (list ForecastList) RecommendCache(repo []*pricing.Repository) (pricing.RecommendedList, error) {
	rmap := make(map[string]*pricing.Repository)
	for i := range repo {
		rmap[repo[i].Region[0]] = repo[i]
	}

	out := pricing.RecommendedList{}
	for _, in := range list {
		if len(in.CacheEngine) < 1 {
			continue
		}

		repo := rmap[in.Region]

		// https://aws.amazon.com/elasticache/reserved-cache-nodes/
		// For latest generation nodes (M5, R5 onwards),
		// you can choose between three payment options
		// when you purchase a Reserved Instance.
		// With the All Upfront option,
		// you pay for the entire Reserved Instance with one upfront payment.
		price := repo.SelectAll().
			UsageType(in.UsageType).
			LeaseContractLength("1yr").
			PurchaseOptionOR([]string{"All Upfront", "Heavy Utilization"}).
			CacheEngine(in.CacheEngine)

		if len(price) != 1 {
			continue
		}

		forecast := in.InstanceNum.ForecastList()
		rec, err := repo.Recommend(price[0], forecast, "minimum")
		if err != nil {
			return nil, fmt.Errorf("recommend cache: %v", err)
		}

		if rec.ReservedInstanceNum > 0 {
			out = append(out, rec)
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Record.UsageType < out[j].Record.UsageType
	})

	return out, nil
}

func (list ForecastList) RecommendDatabase(repo []*pricing.Repository) (pricing.RecommendedList, error) {
	rmap := make(map[string]*pricing.Repository)
	for i := range repo {
		rmap[repo[i].Region[0]] = repo[i]
	}

	out := pricing.RecommendedList{}

	for _, in := range list {
		if len(in.DatabaseEngine) < 1 {
			continue
		}

		repo := rmap[in.Region]
		price := repo.SelectAll().
			UsageType(in.UsageType).
			LeaseContractLength("1yr").
			PurchaseOption("All Upfront").
			DatabaseEngine(in.DatabaseEngine)

		if len(price) != 1 {
			continue
		}

		forecast := in.InstanceNum.ForecastList()
		rec, err := repo.Recommend(price[0], forecast, "minimum")
		if err != nil {
			return nil, fmt.Errorf("recommend database: %v", err)
		}

		if rec.ReservedInstanceNum > 0 {
			out = append(out, rec)
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Record.UsageType < out[j].Record.UsageType
	})

	return out, nil
}

func (list ForecastList) Recommend(repo []*pricing.Repository) (pricing.RecommendedList, error) {
	compute, err := list.RecommendCompute(repo)
	if err != nil {
		return nil, fmt.Errorf("recommend compute: %v", err)
	}

	cache, err := list.RecommendCache(repo)
	if err != nil {
		return nil, fmt.Errorf("recommend cache: %v", err)
	}

	database, err := list.RecommendDatabase(repo)
	if err != nil {
		return nil, fmt.Errorf("recommend database: %v", err)
	}

	out := pricing.RecommendedList{}
	out = append(out, compute...)
	out = append(out, cache...)
	out = append(out, database...)

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Record.UsageType < out[j].Record.UsageType
	})

	return out, nil
}

func (list ForecastList) Region() []string {
	out := []string{}
	for i := range list {
		out = append(out, list[i].Region)
	}

	return out
}

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

func (list ForecastList) Header() []interface{} {
	header := []interface{}{
		"account_id",
		"alias",
		"usage_type",
		"platform/engine",
	}
	for _, n := range list[0].InstanceNum {
		header = append(header, n.Date)
	}

	return header
}

func (list ForecastList) Array() [][]interface{} {
	out := [][]interface{}{}
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

		out = append(out, val)
	}

	return out
}

type InstanceNum struct {
	Date        string  `json:"date"`
	InstanceNum float64 `json:"instance_num"`
}

type InstanceNumList []*InstanceNum

func (list InstanceNumList) ForecastList() pricing.ForecastList {
	forecast := pricing.ForecastList{}
	for _, n := range list {
		forecast = append(forecast, &pricing.Forecast{
			Date:        n.Date,
			InstanceNum: n.InstanceNum,
		})
	}

	return forecast
}

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

type Coverage struct {
	UsageType   string  `json:"usage_type"`
	OSEngine    string  `json:"os_engine"`
	InstanceNum float64 `json:"instance_num"`
	CurrentRI   float64 `json:"current_ri"`
	Short       float64 `json:"short"`
	Coverage    float64 `json:"coverage"`
}

type CoverageList []*Coverage

func (list CoverageList) Header() []interface{} {
	return []interface{}{
		"usage_type",
		"os/engine",
		"instance_num",
		"current_ri",
		"short",
		"coverage",
	}
}

func (list CoverageList) Array() [][]interface{} {
	out := [][]interface{}{}
	for _, r := range list {
		out = append(out, []interface{}{
			r.UsageType,
			r.OSEngine,
			r.InstanceNum,
			r.CurrentRI,
			r.Short,
			r.Coverage,
		})
	}

	return out
}

func GetComputeCoverage(list pricing.NormalizedList, rsv *reserved.Repository) (CoverageList, error) {
	used := reserved.RecordList{}
	out := CoverageList{}
	for i := range list {
		min := list[i].Record
		if !min.IsInstance() {
			continue
		}

		rs := rsv.SelectAll().
			InstanceType(min.InstanceType).
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
			return nil, fmt.Errorf("invalid compute reservation: %v", rs)
		}

		out = append(out, &Coverage{
			UsageType:   min.UsageType,
			OSEngine:    min.OSEngine(),
			InstanceNum: list[i].InstanceNum,
			CurrentRI:   current,
			Short:       list[i].InstanceNum - current,
			Coverage:    current / list[i].InstanceNum,
		})
	}

	unused := reserved.RecordList{}
	for _, r := range rsv.SelectAll().Active() {
		if len(r.InstanceType) < 1 {
			continue
		}

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
			UsageType:   UsageType(r),
			OSEngine:    OSEngine(r),
			InstanceNum: 0,
			CurrentRI:   float64(r.Count()),
			Short:       float64(-r.Count()),
			Coverage:    float64(r.Count()) / 0.0,
		})
	}

	return out, nil
}

func GetCacheCoverage(list pricing.NormalizedList, rsv *reserved.Repository) (CoverageList, error) {
	used := reserved.RecordList{}
	out := CoverageList{}
	for i := range list {
		min := list[i].Record
		if !min.IsCacheNode() {
			continue
		}

		rs := rsv.SelectAll().
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
			InstanceNum: list[i].InstanceNum,
			CurrentRI:   current,
			Short:       list[i].InstanceNum - current,
			Coverage:    current / list[i].InstanceNum,
		})
	}

	unused := reserved.RecordList{}
	for _, r := range rsv.SelectAll().Active() {
		if len(r.CacheNodeType) < 1 {
			continue
		}

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
			UsageType:   UsageType(r),
			OSEngine:    OSEngine(r),
			InstanceNum: 0,
			CurrentRI:   float64(r.Count()),
			Short:       float64(-r.Count()),
			Coverage:    float64(r.Count()) / 0.0,
		})
	}

	return out, nil
}

func GetDatabaseCoverage(list pricing.NormalizedList, rsv *reserved.Repository) (CoverageList, error) {
	used := reserved.RecordList{}
	out := CoverageList{}
	for i := range list {
		min := list[i].Record
		if !min.IsDatabase() {
			continue
		}

		maz := false
		if strings.Contains(min.UsageType, "Multi-AZ") {
			maz = true
		}

		rs := rsv.SelectAll().
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
			InstanceNum: list[i].InstanceNum,
			CurrentRI:   current,
			Short:       list[i].InstanceNum - current,
			Coverage:    current / list[i].InstanceNum,
		})
	}

	unused := reserved.RecordList{}
	for _, r := range rsv.SelectAll().Active() {
		if len(r.DBInstanceClass) < 1 {
			continue
		}

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
			UsageType:   UsageType(r),
			OSEngine:    OSEngine(r),
			InstanceNum: 0,
			CurrentRI:   float64(r.Count()),
			Short:       float64(-r.Count()),
			Coverage:    float64(r.Count()) / 0.0,
		})
	}

	return out, nil
}

func GetCoverage(list pricing.NormalizedList, rsv *reserved.Repository) (CoverageList, error) {
	compute, err := GetComputeCoverage(list, rsv)
	if err != nil {
		return nil, fmt.Errorf("get compute coverage: %v", err)
	}

	cache, err := GetCacheCoverage(list, rsv)
	if err != nil {
		return nil, fmt.Errorf("get cache coverage: %v", err)
	}

	database, err := GetDatabaseCoverage(list, rsv)
	if err != nil {
		return nil, fmt.Errorf("get database coverage: %v", err)
	}

	out := CoverageList{}
	out = append(out, compute...)
	out = append(out, cache...)
	out = append(out, database...)

	return out, nil
}

func UsageType(r *reserved.Record) string {
	region := costexp.Lookup(r.Region)
	if len(r.InstanceType) > 0 {
		return region + "-BoxUsage:" + r.InstanceType
	}

	if len(r.CacheNodeType) > 0 {
		return region + "-NodeUsage:" + r.CacheNodeType
	}

	if len(r.DBInstanceClass) > 0 {
		if r.MultiAZ {
			return region + "-Multi-AZUsage:" + r.DBInstanceClass
		}
		return region + "-InstanceUsage:" + r.DBInstanceClass
	}

	panic("instancetype/cachenodetype/dbinstanceclass not found")
}

func OSEngine(r *reserved.Record) string {
	if len(r.InstanceType) > 0 {
		return pricing.OperatingSystem[r.ProductDescription]
	}

	if len(r.CacheNodeType) > 0 {
		return strings.Title(r.ProductDescription)
	}

	if len(r.DBInstanceClass) > 0 {
		return strings.Replace(strings.Title(r.ProductDescription), "-", " ", -1)
	}

	panic("instancetype/cachenodetype/dbinstanceclass not found")
}
