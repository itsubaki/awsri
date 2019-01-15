package costexp

import (
	"fmt"
	"testing"
)

func TestUnique(t *testing.T) {
	dir := "/var/tmp/hermes/costexp"
	path := fmt.Sprintf("%s/%s", dir, "example_2018-09.out")
	repo, err := NewRepository(path)
	if err != nil {
		t.Errorf("new repository: %v", err)
	}

	for _, r := range repo.SelectAll().Unique("Platform") {
		fmt.Println(r)
	}

	for _, r := range repo.SelectAll().Unique("CacheEngine") {
		fmt.Println(r)
	}

	for _, r := range repo.SelectAll().Unique("DatabaseEngine") {
		fmt.Println(r)
	}
}
