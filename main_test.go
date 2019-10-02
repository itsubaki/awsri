package main

import (
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/pkg/hermes"
	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

// usage forecast by aws account
// normalize usage forecast by date
// merge normalized usage forecast by date
// break-even point with purchase option
func TestPackage(t *testing.T) {
	// price list
	plist, err := pricing.Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize: %v", err)
	}

	// family -> minimum price
	family := pricing.Family(plist)
	for _, v := range family {
		fmt.Println(v)
	}

	mini := pricing.Minimum(family, plist)
	for _, v := range mini {
		fmt.Println(v)
	}

	// forecast quantity
	forecast, err := usage.Deserialize("/var/tmp/hermes", usage.Last12Months())
	if err != nil {
		t.Errorf("usage deserialize: %v", err)
	}

	n := hermes.Normalize(forecast, mini)
	for _, nn := range n {
		fmt.Println(nn)
	}

	merged := usage.Merge(n)
	for _, m := range merged {
		fmt.Println(m)
	}

	monthly := usage.Monthly(merged)
	for _, m := range monthly {
		fmt.Println(m)
	}

	// recommend
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

	for _, pp := range price {
		hash := hermes.Hash(fmt.Sprintf("%s%s", pp.UsageType, "Linux/UNIX"))
		q, p, err := hermes.BreakEvenPoint(monthly[hash], pp)
		if err != nil {
			t.Errorf("%v", err)
		}
		fmt.Printf("%s -> %.0f\n", p, q.InstanceNum)
	}
}
