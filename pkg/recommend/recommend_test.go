package recommend

import (
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/pkg/pricing"

	"github.com/itsubaki/hermes/pkg/usage"
)

func TestRecommend(t *testing.T) {
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
	}

	quantity, err := usage.Deserialize("/var/tmp/hermes", usage.Last12Months())
	if err != nil {
		t.Errorf("usage deserialize: %v", err)
	}

	res, err := Recommend(quantity, price)
	if err != nil {
		t.Errorf("recommend: %v", err)
	}

	for _, r := range res {
		fmt.Printf("%#v\n", r)
	}
}

func TestRecommendFull(t *testing.T) {
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
	}

	quantity, err := usage.Deserialize("/var/tmp/hermes", usage.Last12Months())
	if err != nil {
		t.Errorf("usage deserialize: %v", err)
	}

	res, err := Recommend(quantity, price)
	if err != nil {
		t.Errorf("recommend: %v", err)
	}

	plist, err := pricing.Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize: %v", err)
	}

	normalized, err := Normalize(res, plist)
	if err != nil {
		t.Errorf("normalized: %v", err)
	}

	for _, q := range Merge(normalized) {
		fmt.Printf("%#v\n", q)
	}
}
