package hermes

import (
	"sort"

	"github.com/itsubaki/hermes/pkg/usage"
)

func MonthlyUsage(quantity []usage.Quantity) map[string][]usage.Quantity {
	merged := make(map[string]usage.Quantity)
	for _, q := range quantity {
		hash := q.Hash()
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

	monthly := make(map[string][]usage.Quantity)
	for _, q := range merged {
		hash := q.HashWithOutDate()
		monthly[hash] = append(monthly[hash], q)
	}

	for k := range monthly {
		sort.Slice(monthly[k], func(i, j int) bool { return monthly[k][i].Date < monthly[k][j].Date })
	}

	return monthly
}
