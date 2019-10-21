package reservation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/itsubaki/hermes/pkg/pricing"
)

func Normalize(u []Utilization, mini map[string]pricing.Tuple) []Utilization {
	out := make([]Utilization, 0)
	for i := range u {
		hash := fmt.Sprintf("%s%s", u[i].UsageType(), u[i].OSEngine())
		v, ok := mini[hash]
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
			AccountID:        u[i].AccountID,
			Description:      u[i].Description,
			Region:           u[i].Region,
			InstanceType:     strings.Split(v.Minimum.UsageType, ":")[1],
			Platform:         u[i].Platform,
			CacheEngine:      u[i].CacheEngine,
			DatabaseEngine:   u[i].DatabaseEngine,
			DeploymentOption: u[i].DeploymentOption,
			Date:             u[i].Date,
			Hours:            u[i].Hours * scale,
			Percentage:       u[i].Percentage,
		})
	}

	return out
}
