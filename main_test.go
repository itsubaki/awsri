package main_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/calendar"
	"github.com/itsubaki/hermes/pricing"
	"github.com/itsubaki/hermes/recommend"
	"github.com/itsubaki/hermes/usage"
)

func TestPackage(t *testing.T) {
	price := []pricing.Price{
		{
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
		t.Errorf("desirialize pricing: %v", err)
	}

	family := pricing.Family(plist)
	mini := pricing.Minimum(plist, family)

	date := calendar.LastNMonths(12)
	forecast, err := usage.Deserialize("/var/tmp/hermes", date)
	if err != nil {
		t.Errorf("deserialize usage: %v", err)
	}

	normalized := usage.Normalize(forecast, mini)
	merged := usage.MergeOverall(normalized)
	g, _ := usage.GroupBy(merged)

	for _, p := range price {
		for k := range g {
			if len(g[k][0].Platform) > 0 {
				if p.UsageType != g[k][0].UsageType || p.OperatingSystem != usage.OperatingSystem[g[k][0].Platform] {
					continue
				}
			}

			if len(g[k][0].Platform) < 1 {
				if fmt.Sprintf("%s%s%s", p.UsageType, p.CacheEngine, p.DatabaseEngine) != k {
					continue
				}
			}

			q, _ := recommend.BreakEvenPoint(g[k], p)
			fmt.Println(q)
			break
		}
	}
}
