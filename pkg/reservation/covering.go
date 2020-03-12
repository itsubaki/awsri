package reservation

import (
	"fmt"
	"math"

	"github.com/itsubaki/hermes/pkg/pricing"
)

func AddCoveringCost(plist []pricing.Price, u []Utilization) []string {
	warning := make([]string, 0)

	cache := make(map[string]pricing.Price)
	for i := range plist {
		key := fmt.Sprintf(
			"%s_%s_%s",
			plist[i].UsageType,
			plist[i].OSEngine(),
			plist[i].PreInstalled,
		)

		if v, ok := cache[key]; ok && v.OnDemand != plist[i].OnDemand {
			warning = append(warning, fmt.Sprintf("[WARNING] unexpected pricing: %v", v))
		}

		cache[key] = plist[i]
	}

	for i := range u {
		key := fmt.Sprintf(
			"%s_%s_%s",
			u[i].UsageType(),
			u[i].OSEngine(),
			PreInstalled[u[i].Platform],
		)

		p, ok := cache[key]
		if !ok {
			warning = append(warning, fmt.Sprintf("[WARNING] pricing not found: %v", u[i]))
			continue
		}

		u[i].CoveringCost = math.Round(p.OnDemand*u[i].Hours*1000) / 1000
	}

	return warning
}
