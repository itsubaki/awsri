package recommend

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/itsubaki/hermes/pkg/pricing"
)

func FindMinSize(target pricing.Price, plist []pricing.Price) (pricing.Price, error) {
	tmp := make(map[string]pricing.Price)
	for i := range plist {
		hash := Hash(
			fmt.Sprintf(
				"%s%s%s%s",
				strings.Split(plist[i].UsageType, ".")[0],
				plist[i].OperatingSystem,
				plist[i].CacheEngine,
				plist[i].DatabaseEngine,
			),
		)

		v, ok := tmp[hash]
		if !ok {
			tmp[hash] = plist[i]
			continue
		}

		if strings.LastIndex(v.UsageType, ".") < 1 {
			// instance family is not defined.
			continue
		}

		family0 := v.UsageType[:strings.LastIndex(v.UsageType, ".")]
		family1 := target.UsageType[:strings.LastIndex(target.UsageType, ".")]
		if family0 != family1 {
			// instance family is unmatched.
			// ex. m4.xlarge and m5.large
			continue
		}

		f0, _ := strconv.Atoi(v.NormalizationSizeFactor)
		f1, _ := strconv.Atoi(plist[i].NormalizationSizeFactor)
		if f0 > f1 {
			// tmp[m4.2xlarge] = m4.large
			// tmp[m4.4xlarge] = m4.large
			tmp[hash] = plist[i]
		}
	}

	hash := Hash(
		fmt.Sprintf(
			"%s%s%s%s",
			strings.Split(target.UsageType, ".")[0],
			target.OperatingSystem,
			target.CacheEngine,
			target.DatabaseEngine,
		),
	)

	v, ok := tmp[hash]
	if !ok {
		return pricing.Price{}, fmt.Errorf("pricing not found. target=%#v", target)
	}

	return v, nil
}
