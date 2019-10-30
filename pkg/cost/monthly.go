package cost

import "sort"

func Monthly(c []AccountCost) map[string][]AccountCost {
	monthly := make(map[string][]AccountCost)
	for i := range c {
		monthly[c[i].AccountID] = append(monthly[c[i].AccountID], c[i])
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
