package reservation

import (
	"fmt"

	"github.com/itsubaki/hermes/pkg/pricing"
)

type Cache map[string]pricing.Price

func NewCache(plist []pricing.Price) Cache {
	cache := make(map[string]pricing.Price)

	for i := range plist {
		key := fmt.Sprintf(
			"%s_%s_%s",
			plist[i].UsageType,
			plist[i].OSEngine(),
			plist[i].PreInstalled,
		)

		if v, ok := cache[key]; ok && v.OnDemand != plist[i].OnDemand {
			fmt.Printf("[WARNING] unexpected pricing: %v\n", v)
		}

		cache[key] = plist[i]
	}

	return cache
}

func (c Cache) Find(u Utilization) (pricing.Price, error) {
	key := fmt.Sprintf(
		"%s_%s_%s",
		u.UsageType(),
		u.OSEngine(),
		PreInstalled[u.Platform],
	)

	v, ok := c[key]
	if !ok {
		return pricing.Price{}, fmt.Errorf("pricing not found: %v", u)
	}

	return v, nil
}
