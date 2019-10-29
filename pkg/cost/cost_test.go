package cost

import (
	"fmt"
	"os"
	"testing"

	"github.com/itsubaki/hermes/pkg/usage"
)

func TestFetchCostGroupByLinkedAccount(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	m := usage.LastNMonths(1)[0]
	out, err := FetchCostGroupByLinkedAccount(m.Start, m.End)
	if err != nil {
		t.Errorf("fetch cost group by linked account: %v", err)
	}

	for _, c := range out {
		fmt.Println(c)
	}
}
