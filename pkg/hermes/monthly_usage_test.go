package hermes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/itsubaki/hermes/pkg/usage"
)

func TestMonthlyUsage(t *testing.T) {
	quantity := make([]usage.Quantity, 0)
	date := usage.Last12Months()
	for i := range date {
		file := fmt.Sprintf("/var/tmp/hermes/usage/%s.out", date[i].YYYYMM())
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("file not found: %v", file)
		}

		read, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("read %s: %v", file, err)
		}

		var q []usage.Quantity
		if err := json.Unmarshal(read, &q); err != nil {
			t.Errorf("unmarshal: %v", err)
		}

		quantity = append(quantity, q...)
	}

	monthly := MonthlyUsage(quantity)
	for _, q := range monthly {
		fmt.Printf("%v, ", q[0].UsageType)
		for _, m := range q {
			fmt.Printf("%v, ", m.Date)
		}
		fmt.Println()
	}
}
