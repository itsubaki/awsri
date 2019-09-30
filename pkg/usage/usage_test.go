package usage

import (
	"fmt"
	"os"
	"testing"
)

func TestFetch(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	m := Last12Months()[0]
	list, err := Fetch(m.Start, m.End)
	if err != nil {
		t.Errorf("get usage quantity: %v", err)
	}

	if len(list) < 1 {
		t.Errorf("usage quantity is empty")
	}

	for i := range list {
		fmt.Printf("%#v\n", list[i])
	}
}
