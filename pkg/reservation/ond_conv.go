package reservation

import (
	"fmt"
	"math"

	"github.com/itsubaki/hermes/pkg/pricing"
)

func AddOnDemandConversionCost(plist []pricing.Price, u []Utilization) []error {
	err := make([]error, 0)

	cache := make(map[string]pricing.Price)
	for i := range plist {
		key := fmt.Sprintf(
			"%s_%s_%s",
			plist[i].UsageType,
			plist[i].OSEngine(),
			plist[i].PreInstalled,
		)

		if v, ok := cache[key]; ok && v.OnDemand != plist[i].OnDemand {
			err = append(err, fmt.Errorf("unexpected pricing: %v", v))
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
			err = append(err, fmt.Errorf("pricing not found: %v", u[i]))
			continue
		}

		u[i].OnDemandConversionCost = math.Round(p.OnDemand*u[i].Hours*1000) / 1000
	}

	return err
}
