package usage

import (
	"fmt"
	"os"
	"testing"

	"github.com/itsubaki/hermes/pkg/calendar"
)

func TestFetchLinkedAccount(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	m := calendar.LastNMonths(1)[0]
	list, err := FetchLinkedAccount(m.Start, m.End)
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
