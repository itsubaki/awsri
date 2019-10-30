package recommend

import (
	"testing"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func TestBreakEvenPoint(t *testing.T) {
	price := pricing.Price{
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
	}

	forecast := []usage.Quantity{
		{InstanceNum: 120},
		{InstanceNum: 110},
		{InstanceNum: 100},
		{InstanceNum: 90},
		{InstanceNum: 80},
		{InstanceNum: 70},
		{InstanceNum: 60},
		{InstanceNum: 50},
		{InstanceNum: 40},
		{InstanceNum: 30},
		{InstanceNum: 20},
		{InstanceNum: 10},
	}

	q, _ := BreakEvenPoint(forecast, price)
	if q.InstanceNum != 40 {
		t.Errorf("%v", q.InstanceNum)
	}
}
