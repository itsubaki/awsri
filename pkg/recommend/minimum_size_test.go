package recommend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/itsubaki/hermes/pkg/pricing"
)

func TestFindMinSize(t *testing.T) {
	file := fmt.Sprintf("/var/tmp/hermes/pricing/%s.out", "ap-northeast-1")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Printf("file not found: %v", file)
		os.Exit(1)
	}

	read, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("read %s: %v", file, err)
		os.Exit(1)
	}

	var plist []pricing.Price
	if err := json.Unmarshal(read, &plist); err != nil {
		fmt.Printf("unmarshal: %v", err)
		os.Exit(1)
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
		min, err := FindMinSize(tt, plist)
		if err != nil {
			t.Errorf("find min size: %v", err)
		}

		fmt.Printf("%#v\n", min)
	}
}
