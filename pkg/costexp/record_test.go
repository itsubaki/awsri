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
		"src/github.com/itsubaki/awsri/internal/_serialized/costexp",
		"example_2018-11.out",
	)

	repo, err := NewRepository(path)
	if err != nil {
		t.Errorf("new repository: %v", err)
	}

	for _, r := range repo.SelectAll().Unique("AccountID") {
		fmt.Println(r)
	}

	for _, r := range repo.SelectAll().Unique("UsageType") {
		fmt.Println(r)
	}

	for _, r := range repo.SelectAll().Unique("Platform") {
		fmt.Println(r)
	}

	for _, r := range repo.SelectAll().Unique("Engine") {
		fmt.Println(r)
	}
}
