package cost

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/itsubaki/hermes/calendar"
)

func Serialize(dir string, date []calendar.Date, metrics []string) error {
	path := fmt.Sprintf("%s/cost", dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	for i := range date {
		file := fmt.Sprintf("%s/%s.json", path, date[i].String())
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			continue
		}

		ac, err := Fetch(date[i].Start, date[i].End, metrics...)
		if err != nil {
			return fmt.Errorf("fetch cost (%s, %s): %v", date[i].Start, date[i].End, err)
		}

		b, err := json.Marshal(ac)
		if err != nil {
			return fmt.Errorf("marshal: %v", err)
		}

		if err := os.WriteFile(file, b, os.ModePerm); err != nil {
			return fmt.Errorf("write file: %v", err)
		}

		fmt.Printf("write: %v\n", file)
	}

	return nil
}

func Deserialize(dir string, date []calendar.Date) ([]AccountCost, error) {
	cost := make([]AccountCost, 0)
	for _, d := range date {
		file := fmt.Sprintf("%s/cost/%s.json", dir, d.String())
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return []AccountCost{}, fmt.Errorf("file not found: %v", file)
		}

		read, err := os.ReadFile(file)
		if err != nil {
			return []AccountCost{}, fmt.Errorf("read %s: %v", file, err)
		}

		var u []AccountCost
		if err := json.Unmarshal(read, &u); err != nil {
			return []AccountCost{}, fmt.Errorf("unmarshal: %v", err)
		}

		cost = append(cost, u...)
	}

	sort.SliceStable(cost, func(i, j int) bool { return cost[i].AccountID < cost[j].AccountID })
	sort.SliceStable(cost, func(i, j int) bool { return cost[i].Date < cost[j].Date })

	return cost, nil
}
