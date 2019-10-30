package reservation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/itsubaki/hermes/pkg/usage"
)

func Serialize(dir string, date []usage.Date) error {
	path := fmt.Sprintf("%s/reservation", dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	for i := range date {
		file := fmt.Sprintf("%s/%s.out", path, date[i].YYYYMM())
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			continue
		}

		u, err := Fetch(date[i].Start, date[i].End)
		if err != nil {
			return fmt.Errorf("fetch reservation (%s, %s): %v\n", date[i].Start, date[i].End, err)
		}

		bytes, err := json.Marshal(u)
		if err != nil {
			return fmt.Errorf("marshal: %v\n", err)
		}

		if err := ioutil.WriteFile(file, bytes, os.ModePerm); err != nil {
			return fmt.Errorf("write file: %v\n", err)
		}

		fmt.Printf("write: %v\n", file)
	}

	return nil
}

func Deserialize(dir string, date []usage.Date) ([]Utilization, error) {
	utilization := make([]Utilization, 0)
	for _, d := range date {
		file := fmt.Sprintf("%s/reservation/%s.out", dir, d.YYYYMM())
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return []Utilization{}, fmt.Errorf("file not found: %v", file)
		}

		read, err := ioutil.ReadFile(file)
		if err != nil {
			return []Utilization{}, fmt.Errorf("read %s: %v", file, err)
		}

		var u []Utilization
		if err := json.Unmarshal(read, &u); err != nil {
			return []Utilization{}, fmt.Errorf("unmarshal: %v", err)
		}

		utilization = append(utilization, u...)
	}

	sort.SliceStable(utilization, func(i, j int) bool { return utilization[i].AccountID < utilization[j].AccountID })
	sort.SliceStable(utilization, func(i, j int) bool { return utilization[i].Region < utilization[j].Region })
	sort.SliceStable(utilization, func(i, j int) bool { return utilization[i].Date < utilization[j].Date })

	return utilization, nil
}
