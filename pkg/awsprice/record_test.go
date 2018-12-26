package awsprice

import (
	"fmt"
	"os"
	"testing"
)

func TestBreakevenPoint1yr(t *testing.T) {
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

	rs := repo.FindByInstanceType("m4.large").
		OperatingSystem("Linux").
		Tenancy("Shared").
		PreInstalled("NA").
		OfferingClass("standard").
		LeaseContractLength("1yr")

	for _, r := range rs {
		p := r.BreakevenPointInMonth()

		if r.PurchaseOption == "No Upfront" && p != 1 {
			t.Errorf("invalid breakeven point. purchase=%v, point=%v", r.PurchaseOption, p)
		}
		if r.PurchaseOption == "Partial Upfront" && p != 6 {
			t.Errorf("invalid breakeven point. purchase=%v, point=%v", r.PurchaseOption, p)
		}
		if r.PurchaseOption == "All Upfront" && p != 8 {
			t.Errorf("invalid breakeven point. purchase=%v, point=%v", r.PurchaseOption, p)
		}
	}
}

func TestBreakevenPoint3yr(t *testing.T) {
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

	rs := repo.FindByInstanceType("m4.large").
		OperatingSystem("Linux").
		Tenancy("Shared").
		PreInstalled("NA").
		OfferingClass("standard").
		LeaseContractLength("3yr")

	for _, r := range rs {
		p := r.BreakevenPointInMonth()

		if r.PurchaseOption == "No Upfront" && p != 1 {
			t.Errorf("invalid breakeven point. purchase=%v, point=%v", r.PurchaseOption, p)
		}
		if r.PurchaseOption == "Partial Upfront" && p != 11 {
			t.Errorf("invalid breakeven point. purchase=%v, point=%v", r.PurchaseOption, p)
		}
		if r.PurchaseOption == "All Upfront" && p != 16 {
			t.Errorf("invalid breakeven point. purchase=%v, point=%v", r.PurchaseOption, p)
		}
	}
}

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

func TestExpectNoReserved(t *testing.T) {
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

	forecast := []Forecast{
		{Month: "2018-06", InstanceNum: 10},
		{Month: "2018-07", InstanceNum: 20},
		{Month: "2018-08", InstanceNum: 10},
		{Month: "2018-09", InstanceNum: 20},
		{Month: "2018-10", InstanceNum: 10},
		{Month: "2018-11", InstanceNum: 20},
	}

	n, c := r.ExpectedInstanceNumAndCost(forecast)
	if c.FullOnDemand.Total != c.ReservedApplied.Total {
		t.Errorf("invalid total cost")
	}

	if n.OnDemandInstanceNumAvg != 15 {
		t.Errorf("invalid ondemand instance num")
	}

	if n.ReservedInstanceNum != 0 {
		t.Errorf("invalid reserved instance num")
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

	forecast := []Forecast{
		{Month: "2018-01", InstanceNum: 120.4},
		{Month: "2018-02", InstanceNum: 110.3},
		{Month: "2018-03", InstanceNum: 100.1},
		{Month: "2018-04", InstanceNum: 90.4},
		{Month: "2018-05", InstanceNum: 80.9},
		{Month: "2018-06", InstanceNum: 70.6},
		{Month: "2018-07", InstanceNum: 60.3},
		{Month: "2018-08", InstanceNum: 50.9},
		{Month: "2018-09", InstanceNum: 40.7},
		{Month: "2018-10", InstanceNum: 30.4},
		{Month: "2018-11", InstanceNum: 20.2},
		{Month: "2018-12", InstanceNum: 10.8},
	}

	n, c := r.ExpectedInstanceNumAndCost(forecast)
	if n.OnDemandInstanceNumAvg != 23 {
		t.Errorf("invalid ondemand instance num")
	}

	if n.ReservedInstanceNum != 51 {
		t.Errorf("invalid reserved instance num")
	}

	if c.ReservedQuantity != 36363 {
		t.Errorf("invalid reserved quantity")
	}

	if c.Subtraction < 0 {
		t.Error("invalid substraction")
	}

	if c.DiscountRate < 0 {
		t.Error("invalid discount rate")
	}
}
