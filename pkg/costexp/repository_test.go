package costexp

import (
	"fmt"
	"os"
	"testing"
)

func TestSerialize(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	date := []*Date{
		{Start: "2018-02-01", End: "2018-03-01"},
		{Start: "2018-03-01", End: "2018-04-01"},
		{Start: "2018-04-01", End: "2018-05-01"},
		{Start: "2018-05-01", End: "2018-06-01"},
		{Start: "2018-06-01", End: "2018-07-01"},
		{Start: "2018-07-01", End: "2018-08-01"},
		{Start: "2018-08-01", End: "2018-09-01"},
		{Start: "2018-09-01", End: "2018-10-01"},
		{Start: "2018-10-01", End: "2018-11-01"},
		{Start: "2018-11-01", End: "2018-12-01"},
		{Start: "2018-12-01", End: "2019-01-01"},
		{Start: "2019-01-01", End: "2019-02-01"},
	}

	for i := range date {
		path := fmt.Sprintf("/var/tmp/hermes/costexp/%s.out", date[i].YYYYMM())
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

func TestMergedRepository(t *testing.T) {
	dir := "/var/tmp/hermes/costexp"
	path := []string{
		fmt.Sprintf("%s/%s.out", dir, "2018-03"),
		fmt.Sprintf("%s/%s.out", dir, "2018-04"),
		fmt.Sprintf("%s/%s.out", dir, "2018-05"),
		fmt.Sprintf("%s/%s.out", dir, "2018-06"),
		fmt.Sprintf("%s/%s.out", dir, "2018-07"),
		fmt.Sprintf("%s/%s.out", dir, "2018-08"),
		fmt.Sprintf("%s/%s.out", dir, "2018-09"),
	}

	repo := &Repository{}
	for _, p := range path {
		tmp, err := Read(p)
		if err != nil {
			t.Errorf("read file: %v", err)
		}

		repo.Internal = append(repo.Internal, tmp.Internal...)
	}

	if len(repo.SelectAll()) < 1 {
		t.Errorf("invalid repository")
	}

	for _, ID := range repo.AccountID() {
		if len(ID) != 12 {
			t.Errorf("invalid AWS AccountID")
		}
	}
}

func TestRepository(t *testing.T) {
	path := fmt.Sprintf("/var/tmp/hermes/costexp/%s.out", "2018-09")
	repo, err := Read(path)
	if err != nil {
		t.Errorf("read file: %v", err)
	}

	if len(repo.SelectAll()) < 1 {
		t.Errorf("repository is empty")
	}
}
