package recommend

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func Normalize(quantity []usage.Quantity, price []pricing.Price) ([]usage.Quantity, error) {
	out := make([]usage.Quantity, 0)
	for _, q := range quantity {
		n, err := normalize(q, price)
		if err != nil {
			return []usage.Quantity{}, err
		}

		out = append(out, n)
	}

	sort.SliceStable(out, func(i, j int) bool { return out[i].UsageType < out[j].UsageType })
	sort.SliceStable(out, func(i, j int) bool { return out[i].Platform < out[j].Platform })
	sort.SliceStable(out, func(i, j int) bool { return out[i].CacheEngine < out[j].CacheEngine })
	sort.SliceStable(out, func(i, j int) bool { return out[i].DatabaseEngine < out[j].DatabaseEngine })

	return out, nil
}

func normalize(q usage.Quantity, price []pricing.Price) (usage.Quantity, error) {
	p := make([]pricing.Price, 0)
	for i := range price {
		if q.UsageType != price[i].UsageType {
			continue
		}

		if len(q.Platform) > 0 && OperatingSystem[q.Platform] != price[i].OperatingSystem {
			continue
		}

		if len(q.CacheEngine) > 0 && q.CacheEngine != price[i].CacheEngine {
			continue
		}

		if len(q.DatabaseEngine) > 0 && q.DatabaseEngine != price[i].DatabaseEngine {
			continue
		}

		if !HasFlexibility(price[i]) {
			continue
		}

		p = append(p, price[i])
	}

	if len(p) < 1 {
		// candidate not found.
		return q, nil
	}

	min, err := FindMinimumSize(p[0], price)
	if err != nil {
		return usage.Quantity{}, fmt.Errorf("find minimum size: %v", err)
	}

	f0, _ := strconv.Atoi(p[0].NormalizationSizeFactor)
	f1, _ := strconv.Atoi(min.NormalizationSizeFactor)
	pow := float64(f0) / float64(f1)

	return usage.Quantity{
		Region:         q.Region,
		UsageType:      min.UsageType,
		Platform:       q.Platform,
		DatabaseEngine: q.DatabaseEngine,
		CacheEngine:    q.CacheEngine,
		Date:           q.Date,
		InstanceHour:   q.InstanceHour * pow,
		InstanceNum:    q.InstanceNum * pow,
	}, nil
}

// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/apply_ri.html
// Instance size flexibility does not apply to Reserved Instances
// that are purchased for a specific Availability Zone,
// bare metal instances,
// Reserved Instances with dedicated tenancy,
// and Reserved Instances for Windows,
// Windows with SQL Standard,
// Windows with SQL Server Enterprise,
// Windows with SQL Server Web,
// RHEL, and SLES.
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/apply_ri.html
// Instance size flexibility does not apply to Reserved Instances
// that are purchased for a specific Availability Zone,
// bare metal instances,
// Reserved Instances with dedicated tenancy,
// and Reserved Instances for Windows,
// Windows with SQL Standard,
// Windows with SQL Server Enterprise,
// Windows with SQL Server Web,
// RHEL, and SLES.
func HasFlexibility(p pricing.Price) bool {
	if strings.Contains(p.OperatingSystem, "Windows") {
		return false
	}

	if strings.Contains(p.OperatingSystem, "Red Hat Enterprise Linux") {
		return false
	}

	if strings.Contains(p.OperatingSystem, "SUSE Linux") {
		return false
	}

	if strings.Contains(p.Tenancy, "dedicated") {
		return false
	}

	if strings.Contains(p.InstanceType, "cache") {
		return false
	}

	return true
}
