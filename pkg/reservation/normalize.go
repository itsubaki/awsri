package reservation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func Normalize(u []Utilization, mini map[string]pricing.Tuple) []Utilization {
	out := make([]Utilization, 0)

	for i := range u {
		t := "BoxUsage"
		if len(u[i].CacheEngine) > 0 {
			t = "NodeUsage"
		}
		if len(u[i].DatabaseEngine) > 0 {
			t = "InstanceUsage"
		}

		key := fmt.Sprintf(
			"%s-%s:%s%s",
			region[u[i].Region],
			t,
			u[i].InstanceType,
			fmt.Sprintf(
				"%s%s%s",
				usage.OperatingSystem[u[i].Platform],
				u[i].CacheEngine,
				u[i].DatabaseEngine,
			),
		)

		v, ok := mini[key]
		if !ok {
			out = append(out, u[i])
			continue
		}

		if len(v.Minimum.NormalizationSizeFactor) < 1 {
			out = append(out, u[i])
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

		out = append(out, Utilization{
			AccountID:      u[i].AccountID,
			Region:         u[i].Region,
			InstanceType:   strings.Split(v.Minimum.UsageType, ":")[1],
			Platform:       u[i].Platform,
			CacheEngine:    u[i].CacheEngine,
			DatabaseEngine: u[i].DatabaseEngine,
			Date:           u[i].Date,
			Hours:          u[i].Hours * scale,
		})
	}

	return out
}
