package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/pkg/hermes"
	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func TestPackage(t *testing.T) {
	price := []pricing.Price{
		pricing.Price{
			Region:                  "ap-northeast-1",
			UsageType:               "APN1-BoxUsage:c4.large",
			Tenancy:                 "Shared",
			PreInstalled:            "NA",
			OperatingSystem:         "Linux",
			OfferingClass:           "standard",
			LeaseContractLength:     "1yr",
			PurchaseOption:          "All Upfront",
			OnDemand:                0.126,
			ReservedQuantity:        738,
			ReservedHrs:             0,
			NormalizationSizeFactor: "4",
		},
		//pricing.Price{
		//	Region:                  "ap-northeast-1",
		//	UsageType:               "APN1-InstanceUsage:db.r4.large",
		//	Tenancy:                 "Shared",
		//	DatabaseEngine:          "Aurora MySQL",
		//	OfferingClass:           "standard",
		//	LeaseContractLength:     "1yr",
		//	PurchaseOption:          "All Upfront",
		//	OnDemand:                0.35,
		//	ReservedQuantity:        1704,
		//	ReservedHrs:             0,
		//	NormalizationSizeFactor: "4",
		//},
		//pricing.Price{
		//	Region:              "ap-northeast-1",
		//	UsageType:           "APN1-NodeUsage:cache.r3.large",
		//	Tenancy:             "Shared",
		//	CacheEngine:         "Redis",
		//	OfferingClass:       "standard",
		//	LeaseContractLength: "1yr",
		//	PurchaseOption:      "Heavy Utilization",
		//	OnDemand:            0.273,
		//	ReservedQuantity:    777,
		//	ReservedHrs:         0.089,
		//},
	}

	plist, err := pricing.Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		fmt.Errorf("desirialize pricing: %v", err)
	}

	family := pricing.Family(plist)
	mini := pricing.Minimum(family, plist)

	date := usage.Last12Months()
	forecast, err := usage.Deserialize("/var/tmp/hermes", date)
	if err != nil {
		t.Errorf("deserialize usage: %v", err)
	}

	normalized := hermes.Normalize(forecast, mini)
	merged := usage.MergeOverall(normalized)
	monthly := usage.Monthly(merged)

	for _, p := range price {
		for k := range monthly {
			if len(monthly[k][0].Platform) > 0 {
				os := hermes.OperatingSystem[monthly[k][0].Platform]
				if p.UsageType != monthly[k][0].UsageType || p.OperatingSystem != os {
					continue
				}
			}

			if len(monthly[k][0].Platform) < 1 {
				str := fmt.Sprintf("%s%s%s", p.UsageType, p.CacheEngine, p.DatabaseEngine)
				if str != k {
					continue
				}
			}

			q, p := hermes.BreakEvenPoint(monthly[k], p)

			in := hermes.Purchase{
				Price:    p,
				Quantity: monthly[k],
			}

			bytes, err := json.Marshal(in)
			if err != nil {
				panic(err)
			}

			fmt.Println(string(bytes))

			fmt.Println(q, p)
			break
		}
	}
}
