package costexp

import (
	"fmt"
	"testing"
)

func TestUnique(t *testing.T) {
	path := fmt.Sprintf("/var/tmp/hermes/costexp/%s.out", "2018-09")
	repo, err := Read(path)
	if err != nil {
		t.Errorf("read file: %v", err)
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

	for _, r := range repo.SelectAll().Unique("Region") {
		fmt.Println(r)
	}
}
