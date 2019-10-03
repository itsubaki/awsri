package usage

import (
	"fmt"
	"sort"
)

func Monthly(q []Quantity) map[string][]Quantity {
	monthly := make(map[string][]Quantity)
	for i := range q {
		hash := fmt.Sprintf(
			"%s%s%s%s%s",
			q[i].AccountID,
			q[i].UsageType,
			q[i].Platform,
			q[i].CacheEngine,
			q[i].DatabaseEngine,
		)

		monthly[hash] = append(monthly[hash], q[i])
	}

	for k := range monthly {
		sort.Slice(monthly[k], func(i, j int) bool { return monthly[k][i].Date < monthly[k][j].Date })
	}

	return monthly
}

func SortedKey(monthly map[string][]Quantity) []string {
	keys := make([]string, 0)
	for k := range monthly {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}
