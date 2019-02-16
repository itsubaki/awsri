package costexp

import (
	"fmt"
	"testing"
)

func TestGetCurrentDate(t *testing.T) {
	current := GetCurrentDate()

	for _, c := range current {
		fmt.Printf("%v\n", c)
	}
}
