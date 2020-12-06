package usage_test

import (
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/itsubaki/hermes/pkg/calendar"
	"github.com/itsubaki/hermes/pkg/usage"
)

func TestFetchUsageType(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	merged := make([]string, 0)
	for _, d := range calendar.LastNMonths(1) {
		usageType, err := usage.FetchUsageType(d.Start, d.End)
		if err != nil {
			t.Errorf("get usage type: %v", err)
		}

		merged = append(merged, usageType...)
	}

	unique := make(map[string]bool)
	for i := range merged {
		unique[merged[i]] = true
	}

	sorted := make([]string, 0)
	for k := range unique {
		sorted = append(sorted, k)
	}

	sort.Strings(sorted)

	for _, s := range sorted {
		fmt.Println(s)
	}
}

func TestFetch(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	m := calendar.LastNMonths(1)[0]
	list, err := usage.Fetch(m.Start, m.End)
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

func TestFetchWith(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	m := calendar.LastNMonths(1)[0]
	list, err := usage.FetchWith(m.Start, m.End, []usage.FetchFunc{
		usage.FetchSpotUsage,
	})
	if err != nil {
		t.Errorf("get usage quantity: %v", err)
	}

	for i := range list {
		fmt.Printf("%#v\n", list[i])
	}
}
