package reservation

import (
	"fmt"
	"sort"
)

func GroupBy(u []Utilization) (map[string][]Utilization, []string) {
	group := make(map[string][]Utilization)
	for i := range u {
		hash := fmt.Sprintf(
			"%s%s%s%s%s%s%s",
			u[i].AccountID,
			u[i].Region,
			u[i].InstanceType,
			u[i].Platform,
			u[i].CacheEngine,
			u[i].DatabaseEngine,
			u[i].DeploymentOption,
		)

		group[hash] = append(group[hash], u[i])
	}

	for k := range group {
		sort.Slice(group[k], func(i, j int) bool { return group[k][i].Date < group[k][j].Date })
	}

	return group, SortedKey(group)
}

func SortedKey(group map[string][]Utilization) []string {
	keys := make([]string, 0)
	for k := range group {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}
