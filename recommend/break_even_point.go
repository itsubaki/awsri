package recommend

import (
	"math"
	"sort"

	"github.com/itsubaki/hermes/pricing"
	"github.com/itsubaki/hermes/usage"
)

func BreakEvenPoint(monthly []usage.Quantity, price pricing.Price) (usage.Quantity, pricing.Price) {
	p := price.BreakEvenPoint()
	if len(monthly) < p {
		// dont exceed break-even point
		return usage.Quantity{
			Region:         monthly[0].Region,
			UsageType:      monthly[0].UsageType,
			Platform:       monthly[0].Platform,
			DatabaseEngine: monthly[0].DatabaseEngine,
			CacheEngine:    monthly[0].CacheEngine,
			InstanceNum:    0.0,
		}, price
	}

	num := make([]float64, 0)
	for _, v := range monthly {
		num = append(num, v.InstanceNum)
	}
	sort.SliceStable(num, func(i, j int) bool { return num[i] > num[j] })

	return usage.Quantity{
		Region:         monthly[0].Region,
		UsageType:      monthly[0].UsageType,
		Platform:       monthly[0].Platform,
		DatabaseEngine: monthly[0].DatabaseEngine,
		CacheEngine:    monthly[0].CacheEngine,
		InstanceNum:    math.Floor(num[p-1]),
	}, price
}
