package api

import (
	"fmt"
	"strings"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/reserved"
)

type GetCoverageList func(list pricing.NormalizedList, rsv *reserved.Repository) (CoverageList, error)

func NewGetCoverageList() []GetCoverageList {
	return []GetCoverageList{
		GetComputeCoverageList,
		GetCacheCoverageList,
		GetDatabaseCoverageList,
	}
}

func GetComputeCoverageList(list pricing.NormalizedList, rsv *reserved.Repository) (CoverageList, error) {
	used := reserved.RecordList{}
	out := CoverageList{}
	for i := range list {
		min := list[i].Record
		if !min.Compute() {
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

func GetCacheCoverageList(list pricing.NormalizedList, rsv *reserved.Repository) (CoverageList, error) {
	used := reserved.RecordList{}
	out := CoverageList{}
	for i := range list {
		min := list[i].Record
		if !min.Cache() {
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

func GetDatabaseCoverageList(list pricing.NormalizedList, rsv *reserved.Repository) (CoverageList, error) {
	used := reserved.RecordList{}
	out := CoverageList{}
	for i := range list {
		min := list[i].Record
		if !min.Database() {
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

type RecommendQuery func(repo *pricing.Repository, f *Forecast) pricing.RecordList

func NewRecommendQuery() []RecommendQuery {
	return []RecommendQuery{
		NewComputeRecordList,
		NewCacheRecordList,
		NewDatabaseRecordList,
	}
}

func NewComputeRecordList(repo *pricing.Repository, f *Forecast) pricing.RecordList {
	return repo.SelectAll().
		Compute().
		UsageType(f.UsageType).
		OperatingSystem(pricing.OperatingSystem[f.Platform]).
		LeaseContractLength("1yr").
		PurchaseOption("All Upfront").
		OfferingClass("standard").
		PreInstalled("NA").
		Tenancy("Shared")
}

func NewCacheRecordList(repo *pricing.Repository, f *Forecast) pricing.RecordList {
	return repo.SelectAll().
		Cache().
		UsageType(f.UsageType).
		CacheEngine(f.CacheEngine).
		LeaseContractLength("1yr").
		PurchaseOptionOR([]string{"All Upfront", "Heavy Utilization"})
}

func NewDatabaseRecordList(repo *pricing.Repository, f *Forecast) pricing.RecordList {
	return repo.SelectAll().
		Database().
		UsageType(f.UsageType).
		DatabaseEngine(f.DatabaseEngine).
		LeaseContractLength("1yr").
		PurchaseOption("All Upfront")
}
