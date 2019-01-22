package reserved

import (
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/pkg/awsprice"
)

func TestRepository(t *testing.T) {
	dir := "/var/tmp/hermes/reserved"
	path := fmt.Sprintf("%s/%s", dir, "example.out")
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

func TestGetReserved(t *testing.T) {
	repo, err := NewRepository("/var/tmp/hermes/reserved/example.out")
	if err != nil {
		t.Errorf("new repository: %v", err)
	}
	r := repo.SelectAll()[0]

	{
		path := fmt.Sprintf("%s/%s.out", "/var/tmp/hermes/awsprice", r.Region)
		repo, err := awsprice.NewRepository(path)
		if err != nil {
			t.Errorf("new repository: %v", err)
		}

		price, err := r.Price(repo)
		if err != nil {
			t.Errorf("get price: %v", err)
		}

		fmt.Println(price)
	}
}
