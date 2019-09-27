package recommend

import (
	"fmt"
	"sort"

	"github.com/itsubaki/hermes/pkg/usage"
)

var BREAKEVENPOINT = 8

func Recommend(quantity []usage.Quantity) ([]usage.Quantity, error) {
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

	out := make([]usage.Quantity, 0)
	for _, q := range sorted {
		if len(q) < BREAKEVENPOINT {
			continue
		}

		hrs := make([]float64, 0)
		for _, v := range q {
			hrs = append(hrs, v.InstanceHour)
		}
		sort.Float64s(hrs)

		num := make([]float64, 0)
		for _, v := range q {
			num = append(num, v.InstanceNum)
		}
		sort.Float64s(num)

		out = append(out, usage.Quantity{
			Region:         q[0].Region,
			UsageType:      q[0].UsageType,
			Platform:       q[0].Platform,
			DatabaseEngine: q[0].DatabaseEngine,
			CacheEngine:    q[0].CacheEngine,
			InstanceHour:   hrs[BREAKEVENPOINT-1],
			InstanceNum:    num[BREAKEVENPOINT-1],
		})
	}

	sort.SliceStable(out, func(i, j int) bool { return out[i].Region < out[j].Region })
	sort.SliceStable(out, func(i, j int) bool { return out[i].UsageType < out[j].UsageType })
	sort.SliceStable(out, func(i, j int) bool { return out[i].Platform < out[j].Platform })
	sort.SliceStable(out, func(i, j int) bool { return out[i].CacheEngine < out[j].CacheEngine })
	sort.SliceStable(out, func(i, j int) bool { return out[i].DatabaseEngine < out[j].DatabaseEngine })

	return out, nil
}
