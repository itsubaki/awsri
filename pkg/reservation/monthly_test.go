package reservation

import (
	"fmt"
	"os"
	"testing"
)

func TestMonthly(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	merged := make([]Utilization, 0)
	for _, d := range Last12Months() {
		u, err := Fetch(d.Start, d.End)
		if err != nil {
			t.Errorf("fetch: %v", err)
		}

		merged = append(merged, u...)
	}

	for _, m := range Monthly(merged) {
		fmt.Printf("%v, %v, %v, ", m[0].AccountID, m[0].Region, m[0].InstanceType)
		for _, mm := range m {
			fmt.Printf("%v: %v, ", mm.Date, mm.Hours)
		}
		fmt.Println()
	}
}
