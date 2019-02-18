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

	for _, r := range repo.SelectAll().Unique("CacheEngine") {
		if r != "Redis" && r != "Memcached" {
			t.Errorf("invalid cache engine=%v", r)
		}
	}

	for _, r := range repo.SelectAll().Unique("DatabaseEngine") {
		if r != "MySQL" && r != "Aurora MySQL" && r != "PostgreSQL" && r != "Aurora PostgreSQL" {
			t.Errorf("invalid database engine=%v", r)
		}
	}

}
