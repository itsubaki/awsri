package recommend

import (
	"fmt"
	"os"
	"testing"

	"github.com/itsubaki/hermes/pkg/pricing"
)

func TestBreakEvenPoint(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	fmt.Println(BreakEvenPoint(
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
		}),
	)

	fmt.Println(BreakEvenPoint(
		pricing.Price{
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
		}),
	)

	fmt.Println(BreakEvenPoint(
		pricing.Price{
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
		}),
	)
}
