package pricing

import (
	"fmt"
	"sort"
	"testing"
)

func TestFetchRedshift(t *testing.T) {
	p, err := Fetch(Redshift, "ap-northeast-1")
	if err != nil {
		t.Error(err)
	}

	list := make([]Price, 0)
	for _, v := range p {
		list = append(list, v)
	}
	sort.SliceStable(list, func(i, j int) bool { return list[i].UsageType < list[j].UsageType })

	for i := range list {
		fmt.Println(list[i])
	}
}

func TestBreakEvenPoint(t *testing.T) {
	fmt.Println(
		Price{
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
		}.BreakEvenPoint(),
	)

	fmt.Println(
		Price{
			Region:                  "ap-northeast-1",
			UsageType:               "APN1-BoxUsage:c4.large",
			Tenancy:                 "Shared",
			PreInstalled:            "NA",
			OperatingSystem:         "Linux",
			OfferingClass:           "standard",
			LeaseContractLength:     "1yr",
			PurchaseOption:          "Partial Upfront",
			OnDemand:                0.126,
			ReservedQuantity:        377,
			ReservedHrs:             0.043,
			NormalizationSizeFactor: "4",
		}.BreakEvenPoint(),
	)

	fmt.Println(
		Price{
			Region:                  "ap-northeast-1",
			UsageType:               "APN1-BoxUsage:c4.large",
			Tenancy:                 "Shared",
			PreInstalled:            "NA",
			OperatingSystem:         "Linux",
			OfferingClass:           "standard",
			LeaseContractLength:     "1yr",
			PurchaseOption:          "No Upfront",
			OnDemand:                0.126,
			ReservedQuantity:        0,
			ReservedHrs:             0.09,
			NormalizationSizeFactor: "4",
		}.BreakEvenPoint(),
	)
}
