package cost

import (
	"fmt"
	"sort"
)

func GroupBy(c []AccountCost) (map[string][]AccountCost, []string) {
	group := make(map[string][]AccountCost)
	for i := range c {
		hash := fmt.Sprintf("%s_%s_%s", c[i].AccountID, c[i].Service, c[i].RecordType)
		group[hash] = append(group[hash], c[i])
	}

	for k := range group {
		sort.Slice(group[k], func(i, j int) bool { return group[k][i].Date < group[k][j].Date })
	}

	return group, SortedKey(group)
}

func SortedKey(group map[string][]AccountCost) []string {
	keys := make([]string, 0)
	for k := range group {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}
