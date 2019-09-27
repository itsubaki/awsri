package recommend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func TestRecommend(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	quantity := make([]usage.Quantity, 0)
	date := usage.Last12Months()
	for i := range date {
		file := fmt.Sprintf("/var/tmp/hermes/usage/%s.out", date[i].YYYYMM())
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("file not found: %v", file)
		}

		read, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("read %s: %v", file, err)
		}

		var q []usage.Quantity
		if err := json.Unmarshal(read, &q); err != nil {
			t.Errorf("unmarshal: %v", err)
		}

		quantity = append(quantity, q...)
	}

	monthly := MonthlyUsage(quantity)
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

	recommended := make([]usage.Quantity, 0)
	for _, p := range price {
		res, err := Recommend(monthly, p)
		if err != nil {
			t.Errorf("recommend: %v", err)
		}

		recommended = append(recommended, res)
	}

	for _, r := range recommended {
		fmt.Printf("%#v\n", r)
	}

	normalized := make([]usage.Quantity, 0)
	for _, r := range recommended {
		n, err := Normalize(r, price)
		if err != nil {
			t.Errorf("recommend: %v", err)
		}

		normalized = append(normalized, n)
	}

	for _, r := range normalized {
		fmt.Printf("%#v\n", r)
	}
}
