package costexp

import (
	"fmt"
	"os"
	"testing"
)

func TestRepository(t *testing.T) {
	path := fmt.Sprintf(
		"%s/%s/%s.out",
		os.Getenv("GOPATH"),
		"src/github.com/itsubaki/awsri/internal/_serialized/costexp",
		"example_2018-11",
	)

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
