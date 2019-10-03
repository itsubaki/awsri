package usage

import "sort"

func Monthly(quantity []Quantity) map[string][]Quantity {
	monthly := make(map[string][]Quantity)
	for i := range quantity {
		hash := quantity[i].HashWithoutDate()
		monthly[hash] = append(monthly[hash], quantity[i])
	}

	for k := range monthly {
		sort.Slice(monthly[k], func(i, j int) bool { return monthly[k][i].Date < monthly[k][j].Date })
	}

	return monthly
}
