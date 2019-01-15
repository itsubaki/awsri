package costexp

import (
	"fmt"
	"testing"
)

var dir = "/var/tmp/hermes/costexp"

func TestMergedRepository(t *testing.T) {
	path := []string{
		fmt.Sprintf("%s/%s", dir, "example_2018-01.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-02.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-03.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-04.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-05.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-06.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-07.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-08.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-09.out"),
	}

	repo := &Repository{
		Profile: "example",
	}

	for _, p := range path {
		p, err := NewRepository(p)
		if err != nil {
			t.Errorf("new costexp repository: %v", err)
		}
		repo.Internal = append(repo.Internal, p.Internal...)
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
	path := fmt.Sprintf("%s/%s", dir, "example_2018-09.out")
	repo, err := NewRepository(path)
	if err != nil {
		t.Errorf("new repository: %v", err)
	}

	if len(repo.SelectAll()) < 1 {
		t.Errorf("repository is empty")
	}

	if repo.Profile != "example" {
		t.Errorf("invalid profile")
	}

}
