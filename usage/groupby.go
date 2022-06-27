package usage

import (
	"fmt"
	"sort"
)

func GroupBy(q []Quantity) (map[string][]Quantity, []string) {
	group := make(map[string][]Quantity)
	for i := range q {
		hash := fmt.Sprintf(
			"%s%s%s%s%s",
			q[i].AccountID,
			q[i].UsageType,
			q[i].Platform,
			q[i].CacheEngine,
			q[i].DatabaseEngine,
		)

		group[hash] = append(group[hash], q[i])
	}

	for k := range group {
		sort.Slice(group[k], func(i, j int) bool { return group[k][i].Date < group[k][j].Date })
	}

	return group, SortedKey(group)
}

func SortedKey(group map[string][]Quantity) []string {
	keys := make([]string, 0)
	for k := range group {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}
