package billing

import (
	"fmt"
	"os"
	"testing"
)

func TestSerialize(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	date := []*Date{
		{Start: "2018-05-01", End: "2018-06-01"},
		{Start: "2018-06-01", End: "2018-07-01"},
		{Start: "2018-07-01", End: "2018-08-01"},
		{Start: "2018-08-01", End: "2018-09-01"},
		{Start: "2018-09-01", End: "2018-10-01"},
		{Start: "2018-10-01", End: "2018-11-01"},
		{Start: "2018-11-01", End: "2018-12-01"},
		{Start: "2018-12-01", End: "2019-01-01"},
		{Start: "2019-01-01", End: "2019-02-01"},
		{Start: "2019-02-01", End: "2019-03-01"},
		{Start: "2019-03-01", End: "2019-04-01"},
		{Start: "2019-04-01", End: "2019-05-01"},
	}

	for i := range date {
		path := fmt.Sprintf("/var/tmp/hermes/billing/%s.out", date[i].YYYYMM())
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			continue
		}

		repo := NewRepository([]*Date{date[i]})
		if err := repo.Fetch(); err != nil {
			t.Errorf("new repository: %v", err)
		}

		if err := repo.Write(path); err != nil {
			t.Errorf("write file: %v", err)
		}
	}
}
