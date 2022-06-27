package reserved_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/itsubaki/hermes/reservation/reserved"
)

func TestFetch(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	r, err := reserved.Fetch([]string{
		"ap-northeast-1",
		"us-west-1",
		"us-west-2",
		"us-east-1",
		"us-east-2",
	})
	if err != nil {
		t.Errorf("fetch reserved: %v", err)
	}

	for _, rr := range r {
		fmt.Println(rr)
	}
}
