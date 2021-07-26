package usage

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/itsubaki/hermes/pkg/calendar"
)

func Serialize(dir string, date []calendar.Date) error {
	path := fmt.Sprintf("%s/usage", dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	for i := range date {
		file := fmt.Sprintf("%s/%s.json", path, date[i].String())
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			continue
		}

		u, err := Fetch(date[i].Start, date[i].End)
		if err != nil {
			return fmt.Errorf("fetch usage (%s, %s): %v\n", date[i].Start, date[i].End, err)
		}

		b, err := json.Marshal(u)
		if err != nil {
			return fmt.Errorf("marshal: %v\n", err)
		}

		if err := os.WriteFile(file, b, os.ModePerm); err != nil {
			return fmt.Errorf("write file: %v\n", err)
		}

		fmt.Printf("write: %v\n", file)
	}

	return nil
}

func Deserialize(dir string, date []calendar.Date) ([]Quantity, error) {
	quantity := make([]Quantity, 0)
	for _, d := range date {
		file := fmt.Sprintf("%s/usage/%s.json", dir, d.String())
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return []Quantity{}, fmt.Errorf("file not found: %v", file)
		}

		read, err := os.ReadFile(file)
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
