package recommend

import (
	"fmt"
	"sort"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func Recommend(monthly []usage.Quantity, price pricing.Price) (usage.Quantity, error) {
	p := price.BreakEvenPoint()
	if len(monthly) < p {
		return usage.Quantity{}, fmt.Errorf("dont exceed the break-even point %v < %v", len(monthly), p)
	}

	if monthly[0].UsageType != price.UsageType {
		return usage.Quantity{}, fmt.Errorf("usage type is unmatched")
	}

	if len(monthly[0].Platform) > 0 && OperatingSystem[monthly[0].Platform] != price.OperatingSystem {
		return usage.Quantity{}, fmt.Errorf("platform is unmatched")
	}

	if len(monthly[0].CacheEngine) > 0 && monthly[0].CacheEngine != price.CacheEngine {
		return usage.Quantity{}, fmt.Errorf("cache engine is unmatched")
	}

	if len(monthly[0].DatabaseEngine) > 0 && monthly[0].DatabaseEngine != price.DatabaseEngine {
		return usage.Quantity{}, fmt.Errorf("database engine is unmatched")
	}

	hrs, num := make([]float64, 0), make([]float64, 0)
	for _, v := range monthly {
		hrs, num = append(hrs, v.InstanceHour), append(num, v.InstanceNum)
	}
	sort.Float64s(hrs)
	sort.Float64s(num)

	return usage.Quantity{
		Region:         monthly[0].Region,
		UsageType:      monthly[0].UsageType,
		Platform:       monthly[0].Platform,
		DatabaseEngine: monthly[0].DatabaseEngine,
		CacheEngine:    monthly[0].CacheEngine,
		InstanceHour:   hrs[p-1],
		InstanceNum:    num[p-1],
	}, nil
}
