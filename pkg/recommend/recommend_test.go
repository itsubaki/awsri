package recommend

import (
	"fmt"
	"os"
	"testing"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func TestRecommend(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	quantity, err := usage.Deserialize("/var/tmp/hermes", usage.Last12Months())
	if err != nil {
		t.Errorf("usage deserialize: %v", err)
	}

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

	res, err := Recommend(quantity, price)
	if err != nil {
		t.Errorf("recommend: %v", err)
	}

	for _, r := range res {
		fmt.Printf("%#v\n", r)
	}
}
