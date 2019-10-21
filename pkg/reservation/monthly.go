package reservation

import (
	"fmt"
	"sort"
)

func Monthly(u []Utilization) map[string][]Utilization {
	monthly := make(map[string][]Utilization)
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

		monthly[hash] = append(monthly[hash], u[i])
	}

	for k := range monthly {
		sort.Slice(monthly[k], func(i, j int) bool { return monthly[k][i].Date < monthly[k][j].Date })
	}

	return monthly
}

func SortedKey(monthly map[string][]Utilization) []string {
	keys := make([]string, 0)
	for k := range monthly {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}
