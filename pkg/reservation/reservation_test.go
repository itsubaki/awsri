package reservation

import (
	"fmt"
	"os"
	"testing"
)

func TestReservation(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	u, err := Fetch("2019-09-01", "2019-10-01")
	if err != nil {
		t.Errorf("fetch: %v", err)
	}

	for _, uu := range u {
		fmt.Println(uu)
	}
}
