package recommend

import (
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func TestRecommend(t *testing.T) {
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
		{Date: "2019-12", InstanceNum: 120},
		{Date: "2019-11", InstanceNum: 110},
		{Date: "2019-10", InstanceNum: 100},
		{Date: "2019-09", InstanceNum: 90},
		{Date: "2019-08", InstanceNum: 80},
		{Date: "2019-07", InstanceNum: 70},
		{Date: "2019-06", InstanceNum: 60},
		{Date: "2019-05", InstanceNum: 50},
		{Date: "2019-04", InstanceNum: 40},
		{Date: "2019-03", InstanceNum: 30},
		{Date: "2019-02", InstanceNum: 20},
		{Date: "2019-01", InstanceNum: 10},
	}

	fmt.Println(Do(forecast, price).PrettyJSON())
}
