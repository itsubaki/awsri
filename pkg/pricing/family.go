package pricing

import (
	"fmt"
	"strconv"
	"strings"
)

func Family(plist []Price) map[string]Price {
	mmap := make(map[string]Price)
	for i := range plist {
		if strings.LastIndex(plist[i].UsageType, ".") < 0 {
			continue
		}

		hash := fmt.Sprintf(
			"%s%s%s%s%s%s%s%s%s%s%s%s",
			plist[i].UsageType[:strings.LastIndex(plist[i].UsageType, ".")],
			plist[i].OperatingSystem,
			plist[i].CacheEngine,
			plist[i].DatabaseEngine,
			plist[i].Operation,
			plist[i].PreInstalled,
			plist[i].Region,
			plist[i].OfferingClass,
			plist[i].Tenancy,
			plist[i].PurchaseOption,
			plist[i].LeaseContractLength,
			plist[i].Version,
		)

		v, ok := mmap[hash]
		if !ok {
			mmap[hash] = plist[i]
			continue
		}

		if len(v.NormalizationSizeFactor) < 1 {
			continue
		}

		s0, err := strconv.ParseFloat(v.NormalizationSizeFactor, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid normalization size factor: %v", err))
		}

		s1, err := strconv.ParseFloat(plist[i].NormalizationSizeFactor, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid normalization size factor: %v", err))
		}

		if s0 > s1 {
			mmap[hash] = plist[i]
		}
	}

	return mmap
}
