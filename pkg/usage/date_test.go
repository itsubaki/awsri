package usage

import (
	"fmt"
	"testing"
)

func TestLastNMonths(t *testing.T) {
	for _, d := range LastNMonths(1) {
		fmt.Println(d)
	}
}
