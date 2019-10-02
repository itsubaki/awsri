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
			"%s%s%s%s",
			plist[i].UsageType[:strings.LastIndex(plist[i].UsageType, ".")],
			plist[i].OperatingSystem,
			plist[i].CacheEngine,
			plist[i].DatabaseEngine,
		)

		v, ok := mmap[hash]
		if !ok {
			mmap[hash] = plist[i]
			continue
		}

		if len(v.NormalizationSizeFactor) < 1 {
			continue
		}

		s0, _ := strconv.ParseFloat(v.NormalizationSizeFactor, 64)
		s1, _ := strconv.ParseFloat(plist[i].NormalizationSizeFactor, 64)
		if s0 > s1 {
			mmap[hash] = plist[i]
		}
	}

	return mmap
}
