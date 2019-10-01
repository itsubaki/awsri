package hermes

import (
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func TestBreakEvenPoint(t *testing.T) {
	forecast, err := usage.Deserialize("/var/tmp/hermes", usage.Last12Months())
	if err != nil {
		t.Errorf("usage deserialize: %v", err)
	}
	monthly := MonthlyUsage(forecast)

	p := pricing.Price{
		Region:                  "ap-northeast-1",
		UsageType:               "APN1-BoxUsage:c4.8xlarge",
		Tenancy:                 "Shared",
		PreInstalled:            "NA",
		OperatingSystem:         "Linux",
		OfferingClass:           "standard",
		LeaseContractLength:     "1yr",
		PurchaseOption:          "All Upfront",
		OnDemand:                0.126 * 16,
		ReservedQuantity:        738 * 16,
		ReservedHrs:             0 * 16,
		NormalizationSizeFactor: "64",
	}

	hash := Hash(fmt.Sprintf("%s%s", p.UsageType, "Linux/UNIX"))
	q, p, err := BreakEvenPoint(monthly[hash], p)
	if err != nil {
		t.Errorf("%v", err)
	}
	fmt.Printf("%s %s\n", q, p)
}
