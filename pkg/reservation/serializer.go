package reservation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/itsubaki/hermes/pkg/usage"
)

func Serialize(dir string, quantity []Utilization) error {
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
