package awsprice

import (
	"fmt"
	"os"
	"testing"
)

func TestFindByInstanceType(t *testing.T) {
	path := fmt.Sprintf("%s/%s/%s.out", os.Getenv("GOPATH"), "src/github.com/itsubaki/awsri/internal/_serialized/awsprice", "ap-northeast-1")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("file not found: %v", path)
	}

	repo, err := NewRepository(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	rs := repo.FindByInstanceType("m4.large").
		OperatingSystem("Linux").
		Tenancy("Shared").
		PreInstalled("NA").
		OfferingClass("standard")

	for _, r := range rs {
		if r.InstanceType != "m4.large" {
			t.Error("invalid instance type")
		}
		if r.OperatingSystem != "Linux" {
			t.Error("invalid operationg system")
		}
		if r.Tenancy != "Shared" {
			t.Error("invalid tenancy")
		}
	}

	for _, r := range rs {
		if r.GetAnnualCost().Subtraction < 0 {
			t.Error("invalid subtraction")
		}
		if r.GetAnnualCost().DiscountRate < 0 {
			t.Error("invalid discount rate")
		}
	}
}
