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
			UsageType:           "APN1-BoxUsage:c4.large",
			Tenancy:             "Shared",
			PreInstalled:        "NA",
			OperatingSystem:     "Linux",
			OfferingClass:       "standard",
			LeaseContractLength: "1yr",
			PurchaseOption:      "All Upfront",
		},
		pricing.Price{
			UsageType:           "APN1-BoxUsage:c4.xlarge",
			Tenancy:             "Shared",
			PreInstalled:        "NA",
			OperatingSystem:     "Linux",
			OfferingClass:       "standard",
			LeaseContractLength: "1yr",
			PurchaseOption:      "No Upfront",
		},
		pricing.Price{
			UsageType:           "APN1-BoxUsage:c4.2xlarge",
			Tenancy:             "Shared",
			PreInstalled:        "NA",
			OperatingSystem:     "Linux",
			OfferingClass:       "standard",
			LeaseContractLength: "1yr",
			PurchaseOption:      "Partial Upfront",
		},
		pricing.Price{
			UsageType:           "APN1-BoxUsage:c4.4xlarge",
			Tenancy:             "Shared",
			PreInstalled:        "NA",
			OperatingSystem:     "Linux",
			OfferingClass:       "standard",
			LeaseContractLength: "3yr",
			PurchaseOption:      "All Upfront",
		},
	}

	for _, tt := range target {
		m := FindMinimumSize(tt, price)
		if m.UsageType != "APN1-BoxUsage:c4.large" {
			t.Errorf("usage type: %v", m.UsageType)
		}

		fmt.Println(m)
	}
}
