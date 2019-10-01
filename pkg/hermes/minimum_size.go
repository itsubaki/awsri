package hermes

import (
	"strconv"
	"strings"

	"github.com/itsubaki/hermes/pkg/pricing"
)

func MinimumSize(target pricing.Price, price []pricing.Price) pricing.Price {
	tmp := make(map[string]pricing.Price)
	for i := range price {
		h := price[i].Hash()

		v, ok := tmp[h]
		if !ok {
			tmp[h] = price[i]
			continue
		}

		if strings.LastIndex(v.UsageType, ".") < 1 || strings.LastIndex(target.UsageType, ".") < 1 {
			// instance family is not defined.
			continue
		}

		f0 := v.UsageType[:strings.LastIndex(v.UsageType, ".")]
		f1 := target.UsageType[:strings.LastIndex(target.UsageType, ".")]
		if f0 != f1 {
			// instance family is unmatched.
			continue
		}

		s0, _ := strconv.Atoi(v.NormalizationSizeFactor)
		s1, _ := strconv.Atoi(price[i].NormalizationSizeFactor)
		if s0 > s1 {
			// tmp[m4.2xlarge] = m4.large
			// tmp[m4.4xlarge] = m4.large
			tmp[h] = price[i]
		}
	}

	h := target.Hash()
	v, ok := tmp[h]
	if !ok {
		return target
	}

	return v
}
