package recommend

import (
	"fmt"

	"github.com/itsubaki/hermes/pkg/usage"
)

func Merge(quantity []usage.Quantity) map[string]usage.Quantity {
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

	return merged
}
