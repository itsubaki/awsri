package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/itsubaki/hermes/pkg/hermes"
	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

// usage forecast by aws account
// normalize usage forecast by date
// merge normalized usage forecast by date
// break-even point with purchase option
func TestPackage2(t *testing.T) {
	// price list
	plist, err := pricing.Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize: %v", err)
	}

	// family -> minimum price
	mmap := make(map[string]pricing.Price)
	for i := range plist {
		if strings.LastIndex(plist[i].UsageType, ".") < 0 {
			mmap[hash] = plist[i]
			continue
		}

		hash := fmt.Sprintf(
			"%s%s%s%s",
			plist[i].UsageType[:strings.LastIndex(plist[i].UsageType, ".")],
			plist[i].OperatingSystem,
			plist[i].CacheEngine,
			plist[i].DatabaseEngine,
		)

		v, ok := mmap[hash]
		if !ok {
			mmap[hash] = plist[i]
			continue
		}

		if v.NormalizationSizeFactor == "" || plist[i].NormalizationSizeFactor == "" {
			continue
		}

		s0, _ := strconv.ParseFloat(v.NormalizationSizeFactor, 64)
		s1, _ := strconv.ParseFloat(plist[i].NormalizationSizeFactor, 64)

		if s0 > s1 {
			mmap[hash] = plist[i]
		}
	}

	for _, v := range mmap {
		fmt.Printf(
			"%s, %s, %s\n",
			v.UsageType,
			fmt.Sprintf("%s%s%s", v.OperatingSystem, v.CacheEngine, v.DatabaseEngine),
			v.NormalizationSizeFactor,
		)
	}

	type Tuple struct {
		Price   pricing.Price
		Minimum pricing.Price
	}

	smap := make(map[string]Tuple)
	for i := range plist {
		hash := fmt.Sprintf(
			"%s%s%s%s",
			plist[i].UsageType,
			fmt.Sprintf("%s%s%s",
				plist[i].OperatingSystem,
				plist[i].CacheEngine,
				plist[i].DatabaseEngine,
			),
			plist[i].CacheEngine,
			plist[i].DatabaseEngine,
		)

		if strings.LastIndex(plist[i].UsageType, ".") < 0 {
			smap[hash] = Tuple{plist[i], plist[i]}
			continue
		}

		mhash := fmt.Sprintf(
			"%s%s%s%s",
			plist[i].UsageType[:strings.LastIndex(plist[i].UsageType, ".")],
			plist[i].OperatingSystem,
			plist[i].CacheEngine,
			plist[i].DatabaseEngine,
		)
		smap[hash] = Tuple{plist[i], mmap[mhash]}
	}

	for _, v := range smap {
		fmt.Printf(
			"%s, %s, %s -> %s, %s, %s\n",
			v.Price.UsageType,
			fmt.Sprintf(
				"%s%s%s",
				v.Price.OperatingSystem,
				v.Price.CacheEngine,
				v.Price.DatabaseEngine,
			),
			v.Price.NormalizationSizeFactor,
			v.Minimum.UsageType,
			fmt.Sprintf(
				"%s%s%s",
				v.Minimum.OperatingSystem,
				v.Minimum.CacheEngine,
				v.Minimum.DatabaseEngine,
			),
			v.Minimum.NormalizationSizeFactor,
		)
	}

	// forecast quantity
	forecast, err := usage.Deserialize("/var/tmp/hermes", usage.Last12Months())
	if err != nil {
		t.Errorf("usage deserialize: %v", err)
	}

	n := make([]usage.Quantity, 0)
	for i := range forecast {
		hash := fmt.Sprintf(
			"%s%s%s%s",
			forecast[i].UsageType,
			fmt.Sprintf("%s%s%s",
				hermes.OperatingSystem[forecast[i].Platform],
				forecast[i].CacheEngine,
				forecast[i].DatabaseEngine,
			),
			forecast[i].CacheEngine,
			forecast[i].DatabaseEngine,
		)

		v, ok := smap[hash]
		if !ok {
			n = append(n, forecast[i])
			continue
		}

		if v.Minimum.NormalizationSizeFactor == "" || v.Price.NormalizationSizeFactor == "" {
			continue
		}

		s0, _ := strconv.ParseFloat(v.Minimum.NormalizationSizeFactor, 64)
		s1, _ := strconv.ParseFloat(v.Price.NormalizationSizeFactor, 64)
		scale := s1 / s0

		n = append(n, usage.Quantity{
			AccountID:    forecast[i].AccountID,
			Description:  forecast[i].Description,
			Region:       forecast[i].Region,
			UsageType:    v.Minimum.UsageType,
			Platform:     forecast[i].Platform,
			CacheEngine:  forecast[i].CacheEngine,
			Date:         forecast[i].Date,
			InstanceHour: forecast[i].InstanceHour * scale,
			InstanceNum:  forecast[i].InstanceNum * scale,
		})
	}

	for _, nn := range n {
		fmt.Println(nn)
	}

	merged := make(map[string]usage.Quantity)
	for i := range n {
		v, ok := merged[n[i].Hash()]
		if !ok {
			merged[n[i].Hash()] = usage.Quantity{
				Region:         n[i].Region,
				UsageType:      n[i].UsageType,
				Platform:       n[i].Platform,
				CacheEngine:    n[i].CacheEngine,
				DatabaseEngine: n[i].DatabaseEngine,
				Date:           n[i].Date,
				InstanceHour:   n[i].InstanceHour,
				InstanceNum:    n[i].InstanceNum,
			}
			continue
		}

		merged[n[i].Hash()] = usage.Quantity{
			Region:         v.Region,
			UsageType:      v.UsageType,
			Platform:       v.Platform,
			CacheEngine:    v.CacheEngine,
			DatabaseEngine: v.DatabaseEngine,
			Date:           v.Date,
			InstanceHour:   v.InstanceHour + n[i].InstanceHour,
			InstanceNum:    v.InstanceNum + n[i].InstanceNum,
		}
	}

	for _, n := range merged {
		fmt.Println(n)
	}
}

func TestPackage(t *testing.T) {
	// forecast quantity
	forecast, err := usage.Deserialize("/var/tmp/hermes", usage.Last12Months())
	if err != nil {
		t.Errorf("usage deserialize: %v", err)
	}
	monthly := hermes.MonthlyUsage(forecast)

	// purchase RI
	query := []struct {
		OSEngine string
		Price    pricing.Price
	}{
		{
			OSEngine: "Linux/UNIX",
			Price: pricing.Price{
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
		},
		{
			OSEngine: "Linux/UNIX",
			Price: pricing.Price{
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
			},
		},
		{
			OSEngine: "Linux/UNIX",
			Price: pricing.Price{
				Region:                  "ap-northeast-1",
				UsageType:               "APN1-BoxUsage:c4.2xlarge",
				Tenancy:                 "Shared",
				PreInstalled:            "NA",
				OperatingSystem:         "Linux",
				OfferingClass:           "standard",
				LeaseContractLength:     "1yr",
				PurchaseOption:          "Partial Upfront",
				OnDemand:                0.504,
				ReservedQuantity:        1511,
				ReservedHrs:             0.172,
				NormalizationSizeFactor: "16",
			},
		},
		{
			OSEngine: "Linux/UNIX",
			Price: pricing.Price{
				Region:                  "ap-northeast-1",
				UsageType:               "APN1-BoxUsage:c4.4xlarge",
				Tenancy:                 "Shared",
				PreInstalled:            "NA",
				OperatingSystem:         "Linux",
				OfferingClass:           "standard",
				LeaseContractLength:     "1yr",
				PurchaseOption:          "No Upfront",
				OnDemand:                1.008,
				ReservedQuantity:        0,
				ReservedHrs:             0.722,
				NormalizationSizeFactor: "32",
			},
		},
	}

	plist, err := pricing.Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize: %v", err)
	}

	nn := make([]hermes.Tuple, 0)
	for _, qq := range query {
		hash := hermes.Hash(fmt.Sprintf("%s%s", qq.Price.UsageType, qq.OSEngine))
		q, p, err := hermes.BreakEvenPoint(monthly[hash], qq.Price)
		if err != nil {
			t.Errorf("%v", err)
		}
		//fmt.Printf("%s %s\n", q, p)

		n, p, err := hermes.Normalize(q, p, plist)
		if err != nil {
			t.Errorf("normalize: %v", err)
		}
		//fmt.Printf("%s %s\n", n, p)

		nn = append(nn, hermes.Tuple{Quantity: n, Price: p})
	}

	for _, m := range hermes.Merge(nn) {
		fmt.Println(m)
	}
}
