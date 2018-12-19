package utilization

import (
	"fmt"
	"os"
	"testing"
)

func TestRepository(t *testing.T) {
	awsid := os.Getenv("AWS_ACCOUNT_ID")
	if len(awsid) < 1 {
		return
	}

	path := fmt.Sprintf("%s/%s/%s.out", os.Getenv("GOPATH"), "src/github.com/itsubaki/awsri/internal/_serialized", awsid)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("file not found: %v", path)
	}

	repo, err := NewRepository(path)
	if err != nil {
		t.Errorf("new repository: %v", err)
	}

	for _, r := range repo.SelectAll() {
		fmt.Println(r)
	}

}
