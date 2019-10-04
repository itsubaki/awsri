package main

import (
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

			q, _ := hermes.BreakEvenPoint(monthly[k], p)
			fmt.Println(q)
			break
		}
	}
}
