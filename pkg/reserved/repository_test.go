package reserved

import (
	"os"
	"testing"
)

func TestSerialize(t *testing.T) {
	region := []string{
		"ap-northeast-1",
		"eu-central-1",
		"us-west-1",
		"us-west-2",
	}

	path := "/var/tmp/hermes/reserved/example.out"
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return
	}

	repo, err := NewRepository("example", region)
	if err != nil {
		t.Errorf("new repository: %v", err)
	}

	if err := repo.Write(path); err != nil {
		t.Errorf("write file: %v", err)
	}
}

func TestDeserialize(t *testing.T) {
	repo, err := Read("/var/tmp/hermes/reserved/example.out")
	if err != nil {
		t.Errorf("read file: %v", err)
	}

	if len(repo.SelectAll()) < 1 {
		t.Errorf("repository is empty")
	}

	if repo.Profile != "example" {
		t.Errorf("invalid profile")
	}
}
