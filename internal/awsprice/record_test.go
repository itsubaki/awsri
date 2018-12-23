package awsprice

import (
	"fmt"
	"os"
	"testing"
)

func TestFindByInstanceTypeCache(t *testing.T) {
	path := fmt.Sprintf(
		"%s/%s/%s.out",
		os.Getenv("GOPATH"),
		"src/github.com/itsubaki/awsri/internal/_serialized/awsprice",
		"ap-northeast-1",
	)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("file not found: %v", path)
	}

	repo, err := NewRepository(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	rs := repo.FindByInstanceType("cache.m4.large").
		Engine("Redis").
		PurchaseOption("Heavy Utilization").
		LeaseContractLength("3yr")

	for _, r := range rs {
		if r.Engine != "Redis" {
			t.Error("invalid engine")
		}

		e := r.ExpectedCost(0, 10)
		if e.ReservedApplied.OnDemand != 0 {
			t.Error("invalid reserved applied")
		}

		if e.Subtraction < 0 {
			t.Error("invalid substraction")
		}

		if e.DiscountRate < 0 {
			t.Error("invalid discount rate")
		}
	}
}

func TestExpect(t *testing.T) {
	r := &Record{
		SKU:                     "7MYWT7Y96UT3NJ2D",
		OfferTermCode:           "4NA7Y494T4",
		Region:                  "ap-northeast-1",
		InstanceType:            "m4.large",
		UsageType:               "APN1-BoxUsage:m4.large",
		LeaseContractLength:     "1yr",
		PurchaseOption:          "All Upfront",
		OnDemand:                0.129,
		ReservedHrs:             0,
		ReservedQuantity:        713,
		Tenancy:                 "Shared",
		PreInstalled:            "NA",
		OperatingSystem:         "Linux",
		Operation:               "RunInstances",
		OfferingClass:           "standard",
		NormalizationSizeFactor: "4",
	}

	e := r.ExpectedCost(2, 3)
	if e.FullOnDemand.OnDemand != 5650.2 {
		t.Errorf("invalid full ondemand cost")
	}

	if e.ReservedApplied.OnDemand != 2260.08 {
		t.Errorf("invalid reserved applied cost")
	}

	if e.ReservedApplied.Reserved != 2139 {
		t.Errorf("invalid reserved applied cost")
	}

}
