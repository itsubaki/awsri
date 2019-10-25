package usage

import (
	"fmt"
	"strconv"

	"github.com/itsubaki/hermes/pkg/pricing"
)

func Normalize(q []Quantity, mini map[string]pricing.Tuple) []Quantity {
	out := make([]Quantity, 0)
	for i := range q {
		hash := fmt.Sprintf("%s%s", q[i].UsageType, q[i].OSEngine())
		v, ok := mini[hash]
		if !ok {
			out = append(out, q[i])
			continue
		}

		if len(v.Minimum.NormalizationSizeFactor) < 1 {
			out = append(out, q[i])
			continue
		}

		s0, err := strconv.ParseFloat(v.Minimum.NormalizationSizeFactor, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid normalization size factor: %v", err))
		}

		s1, err := strconv.ParseFloat(v.Price.NormalizationSizeFactor, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid normalization size factor: %v", err))
		}

		scale := s1 / s0

		out = append(out, Quantity{
			AccountID:      q[i].AccountID,
			Description:    q[i].Description,
			Region:         q[i].Region,
			UsageType:      v.Minimum.UsageType,
			Platform:       q[i].Platform,
			CacheEngine:    q[i].CacheEngine,
			DatabaseEngine: q[i].DatabaseEngine,
			Date:           q[i].Date,
			InstanceHour:   q[i].InstanceHour * scale,
			InstanceNum:    q[i].InstanceNum * scale,
			GByte:          q[i].GByte,
			Requests:       q[i].Requests,
			Unit:           q[i].Unit,
		})
	}

	return out
}
