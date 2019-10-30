package account

import (
	"fmt"
	"os"
	"testing"

	"github.com/itsubaki/hermes/pkg/usage"
)

func TestFetchLinkedAccount(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	m := usage.LastNMonths(1)[0]
	list, err := Fetch(m.Start, m.End)
	if err != nil {
		t.Errorf("get usage quantity: %v", err)
	}

	if len(list) < 1 {
		t.Errorf("linked account is empty")
	}

	for i := range list {
		fmt.Printf("%#v\n", list[i])
	}
}
