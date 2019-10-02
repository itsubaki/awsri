package pricing

import (
	"fmt"
	"strings"
)

type Tuple struct {
	Price   Price
	Minimum Price
}

func Minimum(family map[string]Price, plist []Price) map[string]Tuple {
	smap := make(map[string]Tuple)
	for i := range plist {
		hash := fmt.Sprintf(
			"%s%s%s%s",
			plist[i].UsageType,
			plist[i].OperatingSystem,
			plist[i].CacheEngine,
			plist[i].DatabaseEngine,
		)

		if strings.LastIndex(plist[i].UsageType, ".") < 0 {
			smap[hash] = Tuple{plist[i], plist[i]}
			continue
		}

		mhash := fmt.Sprintf(
			"%s%s%s%s",
			plist[i].UsageType[:strings.LastIndex(plist[i].UsageType, ".")],
			plist[i].OperatingSystem,
			plist[i].CacheEngine,
			plist[i].DatabaseEngine,
		)
		smap[hash] = Tuple{plist[i], family[mhash]}
	}

	return smap
}
