package recommend

import (
	"fmt"

	"github.com/itsubaki/hermes/pkg/usage"
)

func MonthlyUsage(quantity []usage.Quantity) map[string][]usage.Quantity {
	merged := make(map[string]usage.Quantity)
	for _, q := range quantity {
		hash := Hash(fmt.Sprintf(
			"%s%s%s%s%s",
			q.UsageType,
			q.Platform,
			q.CacheEngine,
			q.DatabaseEngine,
			q.Date,
		))

		merged[hash] = usage.Quantity{
			Region:         q.Region,
			UsageType:      q.UsageType,
			Platform:       q.Platform,
			DatabaseEngine: q.DatabaseEngine,
			CacheEngine:    q.CacheEngine,
			Date:           q.Date,
			InstanceHour:   merged[hash].InstanceHour + q.InstanceHour,
			InstanceNum:    merged[hash].InstanceNum + q.InstanceNum,
		}
	}

	sorted := make(map[string][]usage.Quantity)
	for _, q := range merged {
		hash := Hash(fmt.Sprintf(
			"%s%s%s%s",
			q.UsageType,
			q.Platform,
			q.CacheEngine,
			q.DatabaseEngine,
		))

		sorted[hash] = append(sorted[hash], q)
	}

	return sorted
}
