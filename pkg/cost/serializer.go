package cost

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/itsubaki/hermes/pkg/usage"
)

func Serialize(dir string, date []usage.Date) error {
	path := fmt.Sprintf("%s/cost", dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	for i := range date {
		file := fmt.Sprintf("%s/%s.out", path, date[i].YYYYMM())
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			continue
		}

		ac, err := Fetch(date[i].Start, date[i].End)
		if err != nil {
			return fmt.Errorf("fetch cost (%s, %s): %v\n", date[i].Start, date[i].End, err)
		}

		b, err := json.Marshal(ac)
		if err != nil {
			return fmt.Errorf("marshal: %v\n", err)
		}

		if err := ioutil.WriteFile(file, b, os.ModePerm); err != nil {
			return fmt.Errorf("write file: %v\n", err)
		}

		fmt.Printf("write: %v\n", file)
	}

	return nil
}

func Deserialize(dir string, date []usage.Date) ([]AccountCost, error) {
	cost := make([]AccountCost, 0)
	for _, d := range date {
		file := fmt.Sprintf("%s/cost/%s.out", dir, d.YYYYMM())
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return []AccountCost{}, fmt.Errorf("file not found: %v", file)
		}

		read, err := ioutil.ReadFile(file)
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
