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

		fhash := fmt.Sprintf(
			//			"%s%s%s%s%s%s%s%s%s%s%s%s",
			"%s%s%s%s%s%s%s%s%s%s",
			plist[i].UsageType[:strings.LastIndex(plist[i].UsageType, ".")],
			plist[i].OperatingSystem,
			plist[i].CacheEngine,
			plist[i].DatabaseEngine,
			//			plist[i].Operation,
			//			plist[i].PreInstalled,
			plist[i].Region,
			plist[i].OfferingClass,
			plist[i].Tenancy,
			plist[i].PurchaseOption,
			plist[i].LeaseContractLength,
			plist[i].Version,
		)

		smap[hash] = Tuple{plist[i], family[fhash]}

		//if strings.Contains(hash, "BoxUsage:c4.8x") && strings.Contains(hash, "Linux") {
		//	if family[fhash].NormalizationSizeFactor == "8" {
		//		fmt.Println(fhash)
		//		fmt.Printf("%s\n", plist[i])
		//	}
		//}
	}

	// validation
	for _, v := range smap {
		if v.Price.OfferTermCode != v.Minimum.OfferTermCode {
			panic("invalid OfferTermCode")
		}
	}

	return smap
}
