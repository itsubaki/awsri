package calendar_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/pkg/calendar"
)

func TestLastNMonths(t *testing.T) {
	for _, d := range calendar.Last12Months() {
		fmt.Println(d)
	}
}
