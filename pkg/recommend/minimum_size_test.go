package recommend

import (
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/pkg/pricing"
)

func TestFindMinimumSize(t *testing.T) {
	price, err := pricing.Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize: %v", err)
	}

	target := []pricing.Price{
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
		pricing.Price{
			Region:                  "ap-northeast-1",
			UsageType:               "APN1-BoxUsage:c4.xlarge",
			Tenancy:                 "Shared",
			PreInstalled:            "NA",
			OperatingSystem:         "Linux",
			OfferingClass:           "standard",
			LeaseContractLength:     "1yr",
			PurchaseOption:          "All Upfront",
			OnDemand:                0.126 * 2,
			ReservedQuantity:        738 * 2,
			ReservedHrs:             0,
			NormalizationSizeFactor: "8",
		},
		pricing.Price{
			Region:                  "ap-northeast-1",
			UsageType:               "APN1-BoxUsage:c4.2xlarge",
			Tenancy:                 "Shared",
			PreInstalled:            "NA",
			OperatingSystem:         "Linux",
			OfferingClass:           "standard",
			LeaseContractLength:     "1yr",
			PurchaseOption:          "All Upfront",
			OnDemand:                0.126 * 4,
			ReservedQuantity:        738 * 4,
			ReservedHrs:             0,
			NormalizationSizeFactor: "16",
		},
		pricing.Price{
			Region:                  "ap-northeast-1",
			UsageType:               "APN1-BoxUsage:c4.4xlarge",
			Tenancy:                 "Shared",
			PreInstalled:            "NA",
			OperatingSystem:         "Linux",
			OfferingClass:           "standard",
			LeaseContractLength:     "1yr",
			PurchaseOption:          "All Upfront",
			OnDemand:                0.126 * 8,
			ReservedQuantity:        738 * 8,
			ReservedHrs:             0,
			NormalizationSizeFactor: "32",
		},
		pricing.Price{
			Region:                  "ap-northeast-1",
			UsageType:               "APN1-NodeUsage:cache.m3.large",
			CacheEngine:             "Redis",
			OfferingClass:           "standard",
			LeaseContractLength:     "3yr",
			PurchaseOption:          "Heavy Utilization",
			OnDemand:                0.24,
			ReservedQuantity:        750,
			ReservedHrs:             0.064,
			NormalizationSizeFactor: "4",
		},
	}

	for _, tt := range target {
		min, err := FindMinimumSize(tt, price)
		if err != nil {
			t.Errorf("find minimum size: %v", err)
		}

		if min.UsageType != "APN1-BoxUsage:c4.large" && min.UsageType != "APN1-NodeUsage:cache.m5.12xlarge" {
			t.Errorf("usage type: %v", min.UsageType)
		}

		fmt.Printf("%#v\n", min)
	}
}
