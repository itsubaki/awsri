package usage

import (
	"fmt"
	"testing"
)

func TestLastNMonths(t *testing.T) {
	for _, d := range Last12Months() {
		fmt.Println(d)
	}
}
