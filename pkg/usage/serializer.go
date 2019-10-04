package usage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

func Serialize(dir string, quantity []Quantity) error {
	return nil
}

func Deserialize(dir string, date []Date) ([]Quantity, error) {
	quantity := make([]Quantity, 0)
	for _, d := range date {
		file := fmt.Sprintf("%s/usage/%s.out", dir, d.YYYYMM())
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return []Quantity{}, fmt.Errorf("file not found: %v", file)
		}

		read, err := ioutil.ReadFile(file)
		if err != nil {
			return []Quantity{}, fmt.Errorf("read %s: %v", file, err)
		}

		var q []Quantity
		if err := json.Unmarshal(read, &q); err != nil {
			return []Quantity{}, fmt.Errorf("unmarshal: %v", err)
		}

		quantity = append(quantity, q...)
	}

	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].AccountID < quantity[j].AccountID })
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].Region < quantity[j].Region })
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].UsageType < quantity[j].UsageType })
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].Date < quantity[j].Date })

	return quantity, nil
}
