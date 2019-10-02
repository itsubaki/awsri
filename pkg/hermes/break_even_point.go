package hermes

import (
	"fmt"
	"math"
	"sort"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func BreakEvenPoint(monthly []usage.Quantity, price pricing.Price) (usage.Quantity, pricing.Price, error) {
	p := price.BreakEvenPoint()
	if len(monthly) < p {
		return usage.Quantity{}, price, fmt.Errorf("dont exceed the break-even point %v < %v", len(monthly), p)
	}

	hrs, num := make([]float64, 0), make([]float64, 0)
	for _, v := range monthly {
		hrs, num = append(hrs, v.InstanceHour), append(num, v.InstanceNum)
	}
	sort.SliceStable(hrs, func(i, j int) bool { return hrs[i] > hrs[j] })
	sort.SliceStable(num, func(i, j int) bool { return num[i] > num[j] })

	return usage.Quantity{
		Region:         monthly[0].Region,
		UsageType:      monthly[0].UsageType,
		Platform:       monthly[0].Platform,
		DatabaseEngine: monthly[0].DatabaseEngine,
		CacheEngine:    monthly[0].CacheEngine,
		InstanceNum:    math.Floor(num[p-1]),
	}, price, nil
}
