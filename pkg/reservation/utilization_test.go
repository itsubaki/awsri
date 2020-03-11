package reservation

import (
	"fmt"
	"os"
	"testing"
)

func TestFetch(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	u, err := Fetch("2019-09-01", "2019-09-02")
	if err != nil {
		t.Errorf("fetch: %v", err)
	}

	for _, uu := range u {
		fmt.Println(uu)
	}
}
