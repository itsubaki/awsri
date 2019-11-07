package cost

import (
	"fmt"
	"os"
	"testing"

	"github.com/itsubaki/hermes/pkg/usage"
)

func TestFetch(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	m := usage.Last12Months()[0]
	out, err := Fetch(m.Start, m.End)
	if err != nil {
		t.Errorf("fetch cost group by linked account: %v", err)
	}

	for _, c := range out {
		fmt.Println(c)
	}
}

func TestFetchWith(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	m := usage.Last12Months()[0]
	out, err := FetchWith(m.Start, m.End, []string{
		"Amazon Elastic Compute Cloud - Compute",
		"Amazon ElastiCache",
		"Amazon Relational Database Service",
	})
	if err != nil {
		t.Errorf("fetch cost group by linked account: %v", err)
	}

	for _, c := range out {
		fmt.Println(c)
	}
}
