package hermes

import (
	"encoding/json"
	"sort"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

type Tuple struct {
	Quantity usage.Quantity `json:"quantity"`
	Price    pricing.Price  `json:"price"`
}

func (t Tuple) String() string {
	return t.JSON()
}

func (t Tuple) JSON() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func Merge(tuple []Tuple) []Tuple {
	tmp := make(map[string]Tuple)
	for _, t := range tuple {
		hash := t.Price.Hash()
		v, ok := tmp[hash]
		if !ok {
			tmp[hash] = t
			continue
		}

		tmp[hash] = Tuple{
			Quantity: usage.Quantity{
				Region:         v.Quantity.Region,
				UsageType:      v.Quantity.UsageType,
				Platform:       v.Quantity.Platform,
				DatabaseEngine: v.Quantity.DatabaseEngine,
				CacheEngine:    v.Quantity.CacheEngine,
				InstanceHour:   v.Quantity.InstanceHour + t.Quantity.InstanceHour,
				InstanceNum:    v.Quantity.InstanceNum + t.Quantity.InstanceNum,
			},
			Price: t.Price,
		}
	}

	out := make([]Tuple, 0)
	for k := range tmp {
		out = append(out, tmp[k])
	}

	sort.SliceStable(out, func(i, j int) bool { return out[i].Quantity.UsageType < out[j].Quantity.UsageType })
	sort.SliceStable(out, func(i, j int) bool { return out[i].Quantity.Platform < out[j].Quantity.Platform })
	sort.SliceStable(out, func(i, j int) bool { return out[i].Quantity.CacheEngine < out[j].Quantity.CacheEngine })
	sort.SliceStable(out, func(i, j int) bool { return out[i].Quantity.DatabaseEngine < out[j].Quantity.DatabaseEngine })

	return out
}
