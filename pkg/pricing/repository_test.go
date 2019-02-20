package pricing

import (
	"fmt"
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

	for i := range region {
		path := fmt.Sprintf("/var/tmp/hermes/pricing/%s.out", region[i])
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			continue
		}

		repo := NewRepository([]string{region[i]})
		if err := repo.Fetch(); err != nil {
			t.Errorf("new repository: %v", err)
		}

		if err := repo.Write(path); err != nil {
			t.Errorf("write file: %v", err)
		}
	}
}

func TestFindByUsageType(t *testing.T) {
	repo, err := Read("/var/tmp/hermes/pricing/ap-northeast-1.out")
	if err != nil {
		t.Errorf("%v", err)
	}

	rs := repo.SelectAll().
		UsageType("APN1-BoxUsage:m4.4xlarge").
		OperatingSystem("Linux").
		Tenancy("Shared").
		LeaseContractLength("1yr").
		PurchaseOption("All Upfront").
		OfferingClass("standard").
		PreInstalled("NA")

	if len(rs) != 1 {
		t.Errorf("invalid result set=%v", rs)
	}
}

func TestNormalizeDatabaseT2Medium(t *testing.T) {
	path := "/var/tmp/hermes/pricing/ap-northeast-1.out"
	repo, err := Read(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	rs := repo.SelectAll().
		InstanceType("db.t2.medium").
		PurchaseOption("All Upfront").
		LeaseContractLength("1yr").
		DatabaseEngine("Aurora MySQL")

	min, err := repo.Normalize(rs[0])
	if err != nil {
		t.Errorf("%v", err)
	}

	if min.UsageType != "APN1-InstanceUsage:db.t2.small" {
		t.Errorf("invalid usage type=%s", min.UsageType)
	}
}

func TestNormalizeDatabase(t *testing.T) {
	path := "/var/tmp/hermes/pricing/ap-northeast-1.out"
	repo, err := Read(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	rs := repo.SelectAll().
		InstanceType("db.m4.4xlarge").
		PurchaseOption("All Upfront").
		LeaseContractLength("1yr").
		DatabaseEngine("PostgreSQL")

	r, err := repo.Normalize(rs[0])
	if err != nil {
		t.Errorf("find normalized record: %v", err)
	}

	if r.InstanceType != "db.m4.large" {
		t.Errorf("invalid instance type=%s in normalized record", r.InstanceType)
	}
}

func TestNormalizeCompute(t *testing.T) {
	path := "/var/tmp/hermes/pricing/ap-northeast-1.out"
	repo, err := Read(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	r := &Record{
		SKU:                     "XU2NYYPCRTK4T7CN",
		OfferTermCode:           "6QCMYABX3D",
		Region:                  "ap-northeast-1",
		InstanceType:            "m4.4xlarge",
		UsageType:               "APN1-BoxUsage:m4.4xlarge",
		LeaseContractLength:     "1yr",
		PurchaseOption:          "All Upfront",
		OnDemand:                1.032,
		ReservedQuantity:        5700,
		ReservedHrs:             0,
		Tenancy:                 "Shared",
		PreInstalled:            "NA",
		OperatingSystem:         "Linux",
		Operation:               "RunInstances",
		OfferingClass:           "standard",
		NormalizationSizeFactor: "32",
	}

	min, err := repo.Normalize(r)
	if err != nil {
		t.Errorf("find normalized record: %v", err)
	}

	if min.InstanceType != "m4.large" {
		t.Errorf("invalid instance type=%s in normalized record", min.InstanceType)
	}
}

func TestFindByInstanceType(t *testing.T) {
	path := "/var/tmp/hermes/pricing/ap-northeast-1.out"
	repo, err := Read(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	rs := repo.SelectAll().
		InstanceType("m4.large").
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
		if r.GetAnnualCost().SavingCost < 0 {
			t.Error("invalid saving cost")
		}
		if r.GetAnnualCost().DiscountRate < 0 {
			t.Error("invalid discount rate")
		}
	}
}

func TestRecommendM4large(t *testing.T) {
	path := "/var/tmp/hermes/pricing/ap-northeast-1.out"
	repo, err := Read(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	r := &Record{
		SKU:                     "7MYWT7Y96UT3NJ2D",
		OfferTermCode:           "6QCMYABX3D",
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

	forecast := ForecastList{
		{Date: "2018-01", InstanceNum: 120.4},
		{Date: "2018-02", InstanceNum: 110.3},
		{Date: "2018-03", InstanceNum: 100.1},
		{Date: "2018-04", InstanceNum: 90.9},
		{Date: "2018-05", InstanceNum: 80.9},
		{Date: "2018-06", InstanceNum: 70.6},
		{Date: "2018-07", InstanceNum: 60.3},
		{Date: "2018-08", InstanceNum: 50.9},
		{Date: "2018-09", InstanceNum: 40.7},
		{Date: "2018-10", InstanceNum: 30.6},
		{Date: "2018-11", InstanceNum: 20.2},
		{Date: "2018-12", InstanceNum: 10.8},
	}

	r0 := r.Recommend(forecast)
	r1, _ := repo.Recommend(r, forecast)

	if r0.Record.UsageType != r1.NormalizedRecord.UsageType {
		t.Errorf("invalid usage type in recommend")
	}
}
