package recommend

import (
	"fmt"
	"sort"

	"github.com/itsubaki/hermes/pkg/usage"
)

func Merge(quantity []usage.Quantity) []usage.Quantity {
	merged := make(map[string]usage.Quantity)
	for i := range quantity {
		hash := Hash(
			fmt.Sprintf(
				"%s%s%s%s",
				quantity[i].UsageType,
				quantity[i].Platform,
				quantity[i].CacheEngine,
				quantity[i].DatabaseEngine,
			),
		)

		v, ok := merged[hash]
		if !ok {
			merged[hash] = usage.Quantity{
				Region:         quantity[i].Region,
				UsageType:      quantity[i].UsageType,
				Platform:       quantity[i].Platform,
				DatabaseEngine: quantity[i].DatabaseEngine,
				CacheEngine:    quantity[i].CacheEngine,
				InstanceHour:   quantity[i].InstanceHour,
				InstanceNum:    quantity[i].InstanceNum,
			}

			continue
		}

		merged[hash] = usage.Quantity{
			Region:         quantity[i].Region,
			UsageType:      quantity[i].UsageType,
			Platform:       quantity[i].Platform,
			DatabaseEngine: quantity[i].DatabaseEngine,
			CacheEngine:    quantity[i].CacheEngine,
			InstanceHour:   quantity[i].InstanceHour + v.InstanceHour,
			InstanceNum:    quantity[i].InstanceNum + v.InstanceNum,
		}
	}

	out := make([]usage.Quantity, 0)
	for k := range merged {
		out = append(out, merged[k])
	}

	sort.SliceStable(out, func(i, j int) bool { return out[i].UsageType < out[j].UsageType })
	sort.SliceStable(out, func(i, j int) bool { return out[i].Platform < out[j].Platform })
	sort.SliceStable(out, func(i, j int) bool { return out[i].CacheEngine < out[j].CacheEngine })
	sort.SliceStable(out, func(i, j int) bool { return out[i].DatabaseEngine < out[j].DatabaseEngine })

	return out
}
