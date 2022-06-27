package reservation_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/itsubaki/hermes/calendar"
	"github.com/itsubaki/hermes/reservation"
)

func TestGroupBy(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	merged := make([]reservation.Utilization, 0)
	for _, d := range calendar.LastNMonths(3) {
		u, err := reservation.Fetch(d.Start, d.End)
		if err != nil {
			t.Errorf("fetch: %v", err)
		}

		merged = append(merged, u...)
	}

	g, _ := reservation.GroupBy(merged)
	for _, m := range g {
		fmt.Printf("%v, %v, %v, ", m[0].AccountID, m[0].Region, m[0].InstanceType)
		for _, mm := range m {
			fmt.Printf("%v: %v, ", mm.Date, mm.Hours)
		}
		fmt.Println()
	}
}
