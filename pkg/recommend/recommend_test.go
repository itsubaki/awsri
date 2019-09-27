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

	// TestData
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

	// Test
	monthly := MonthlyUsage(quantity)
	recommended := make([]usage.Quantity, 0)
	for _, p := range price {
		res, err := Recommend(monthly, p)
		if err != nil {
			// t.Errorf("recommend: %v", err)
			continue
		}

		recommended = append(recommended, res)
	}

	fmt.Println("[recommended]--------------")
	for _, r := range recommended {
		fmt.Printf("%#v\n", r)
	}

	normalized := make([]usage.Quantity, 0)
	for _, r := range recommended {
		n, err := Normalize(r, price)
		if err != nil {
			t.Errorf("normalized: %v", err)
			panic("")
		}

		normalized = append(normalized, n)
	}

	fmt.Println("[normalized]--------------")
	for _, r := range normalized {
		fmt.Printf("%#v\n", r)
	}

	merged := make(map[string]usage.Quantity)
	for i := range normalized {
		hash := Hash(
			fmt.Sprintf(
				"%s%s%s%s",
				normalized[i].UsageType,
				normalized[i].Platform,
				normalized[i].CacheEngine,
				normalized[i].DatabaseEngine,
			),
		)

		v, ok := merged[hash]
		if !ok {
			merged[hash] = usage.Quantity{
				Region:         normalized[i].Region,
				UsageType:      normalized[i].UsageType,
				Platform:       normalized[i].Platform,
				DatabaseEngine: normalized[i].DatabaseEngine,
				CacheEngine:    normalized[i].CacheEngine,
				InstanceHour:   normalized[i].InstanceHour,
				InstanceNum:    normalized[i].InstanceNum,
			}

			continue
		}

		merged[hash] = usage.Quantity{
			Region:         normalized[i].Region,
			UsageType:      normalized[i].UsageType,
			Platform:       normalized[i].Platform,
			DatabaseEngine: normalized[i].DatabaseEngine,
			CacheEngine:    normalized[i].CacheEngine,
			InstanceHour:   normalized[i].InstanceHour + v.InstanceHour,
			InstanceNum:    normalized[i].InstanceNum + v.InstanceNum,
		}
	}

	fmt.Println("[merged]--------------")
	for _, r := range merged {
		fmt.Printf("%#v\n", r)
	}
}
