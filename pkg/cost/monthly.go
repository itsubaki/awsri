package cost

import (
	"fmt"
	"sort"
)

func Monthly(c []AccountCost) map[string][]AccountCost {
	monthly := make(map[string][]AccountCost)
	for i := range c {
		hash := fmt.Sprintf("%s_%s_%s", c[i].AccountID, c[i].Service, c[i].RecordType)
		monthly[hash] = append(monthly[hash], c[i])
	}

	for k := range monthly {
		sort.Slice(monthly[k], func(i, j int) bool { return monthly[k][i].Date < monthly[k][j].Date })
	}

	return monthly
}

func SortedKey(monthly map[string][]AccountCost) []string {
	keys := make([]string, 0)
	for k := range monthly {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}
