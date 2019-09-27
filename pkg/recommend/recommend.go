package recommend

import (
	"fmt"
	"sort"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func Recommend(monthly map[string][]usage.Quantity, price pricing.Price) (usage.Quantity, error) {
	out := make([]usage.Quantity, 0)
	for _, q := range monthly {
		if q[0].UsageType != price.UsageType {
			continue
		}

		if len(q[0].Platform) > 0 && OperatingSystem[q[0].Platform] != price.OperatingSystem {
			continue
		}

		if len(q[0].CacheEngine) > 0 && q[0].CacheEngine != price.CacheEngine {
			continue
		}

		if len(q[0].DatabaseEngine) > 0 && q[0].DatabaseEngine != price.DatabaseEngine {
			continue
		}

		point := price.BreakEvenPoint()
		if len(q) < point {
			continue
		}

		hrs, num := make([]float64, 0), make([]float64, 0)
		for _, v := range q {
			hrs, num = append(hrs, v.InstanceHour), append(num, v.InstanceNum)
		}
		sort.Float64s(hrs)
		sort.Float64s(num)

		out = append(out, usage.Quantity{
			Region:         q[0].Region,
			UsageType:      q[0].UsageType,
			Platform:       q[0].Platform,
			DatabaseEngine: q[0].DatabaseEngine,
			CacheEngine:    q[0].CacheEngine,
			InstanceHour:   hrs[point-1],
			InstanceNum:    num[point-1],
		})
	}

	if len(out) > 1 {
		return usage.Quantity{}, fmt.Errorf("duplicated result. usage=%#v", out)
	}

	if len(out) == 0 {
		return usage.Quantity{}, fmt.Errorf("usage not found. price=%v", price)
	}

	return out[0], nil
}
