package reservation_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/itsubaki/hermes/reservation"
)

func TestFetch(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	u, err := reservation.Fetch("2020-02-01", "2020-02-02")
	if err != nil {
		t.Errorf("fetch: %v", err)
	}

	for _, uu := range u {
		fmt.Println(uu)
	}
}
