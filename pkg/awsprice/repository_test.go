package awsprice

import (
	"fmt"
	"os"
	"testing"
)

func TestFindMinimumDatabase(t *testing.T) {
	path := fmt.Sprintf(
		"%s/%s/%s",
		os.Getenv("GOPATH"),
		"src/github.com/itsubaki/awsri/internal/_serialized/awsprice",
		"ap-northeast-1.out",
	)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("file not found: %v", path)
	}

	repo, err := NewRepository(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	rs := repo.FindByInstanceType("db.m4.4xlarge").
		PurchaseOption("All Upfront").
		LeaseContractLength("1yr").
		DatabaseEngine("PostgreSQL")

	r, err := repo.FindMinimumInstanceType(rs[0])
	if err != nil {
		t.Errorf("find minimum instance type: %v", err)
	}

	if r.InstanceType != "db.m4.large" {
		t.Errorf("invalid minimum instance type=%s", r.InstanceType)
	}
}

func TestFindMinimumCompute(t *testing.T) {
	path := fmt.Sprintf(
		"%s/%s/%s",
		os.Getenv("GOPATH"),
		"src/github.com/itsubaki/awsri/internal/_serialized/awsprice",
		"ap-northeast-1.out",
	)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("file not found: %v", path)
	}

	repo, err := NewRepository(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	rs := repo.FindByInstanceType("m4.4xlarge").
		OperatingSystem("Linux").
		Tenancy("Shared").
		PreInstalled("NA").
		OfferingClass("standard").
		LeaseContractLength("1yr").
		PurchaseOption("All Upfront")

	r, err := repo.FindMinimumInstanceType(rs[0])
	if err != nil {
		t.Errorf("find minimum instance type: %v", err)
	}

	if r.InstanceType != "m4.large" {
		t.Errorf("invalid minimum instance type=%s", r.InstanceType)
	}
}

func TestFindByInstanceType(t *testing.T) {
	path := fmt.Sprintf(
		"%s/%s/%s",
		os.Getenv("GOPATH"),
		"src/github.com/itsubaki/awsri/internal/_serialized/awsprice",
		"ap-northeast-1.out",
	)

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
