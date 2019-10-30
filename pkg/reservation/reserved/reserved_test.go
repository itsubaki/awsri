package reserved

import (
	"fmt"
	"os"
	"testing"
)

func TestFetch(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")
	r, err := Fetch([]string{"ap-northeast-1", "us-west-2"})
	if err != nil {
		t.Errorf("fetch reserved: %v", err)
	}

	for _, rr := range r {
		fmt.Println(rr)
	}
}
