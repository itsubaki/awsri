package hermes

import (
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/pkg/usage"
)

func TestMonthlyUsage(t *testing.T) {
	quantity, err := usage.Deserialize("/var/tmp/hermes", usage.Last12Months())
	if err != nil {
		t.Errorf("usage deserialize: %v", err)
	}

	monthly := MonthlyUsage(quantity)
	for _, q := range monthly {
		fmt.Printf("%v, %s%s%s, ", q[0].UsageType, q[0].Platform, q[0].CacheEngine, q[0].DatabaseEngine)
		for _, m := range q {
			fmt.Printf("%v, ", m.Date)
		}
		fmt.Println()
	}
}
