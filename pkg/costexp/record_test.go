package costexp

import (
	"fmt"
	"os"
	"testing"
)

func TestUnique(t *testing.T) {
	path := fmt.Sprintf(
		"%s/%s/%s",
		os.Getenv("GOPATH"),
		"src/github.com/itsubaki/hermes/internal/_serialized/costexp",
		"example_2018-09.out",
	)

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
