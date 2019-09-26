package usage

import (
	"fmt"
	"os"
	"testing"
)

func TestFetchQuantity(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	list, err := Fetch("2018-11-01", "2018-11-02")
	if err != nil {
		t.Errorf("get usage quantity: %v", err)
	}

	if len(list) < 1 {
		t.Errorf("usage quantity is empty")
	}

	for i := range list {
		fmt.Printf("%v\n", list[i])
	}
}
