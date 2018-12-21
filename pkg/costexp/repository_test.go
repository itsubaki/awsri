package costexp

import (
	"fmt"
	"os"
	"testing"
)

func TestRepository(t *testing.T) {
	path := fmt.Sprintf("%s/%s/%s.out", os.Getenv("GOPATH"), "src/github.com/itsubaki/awsri/internal/_serialized/costexp", "example")

	repo, err := NewRepository(path)
	if err != nil {
		t.Errorf("new repository: %v", err)
	}

	for _, r := range repo.SelectAll() {
		fmt.Println(r)
	}

	if repo.Profile != "example" {
		t.Errorf("invalid profile=%s", repo.Profile)
	}

	if repo.Period.Start != "2018-11-01" {
		t.Errorf("invalid start=%s", repo.Period.Start)
	}

	if repo.Period.End != "2018-12-01" {
		t.Errorf("invalid start=%s", repo.Period.End)
	}
}
