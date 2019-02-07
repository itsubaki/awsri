package reservation

import (
	"os"
	"testing"

	"github.com/itsubaki/hermes/pkg/pricing"
)

func TestSerialize(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")
	region := []string{
		"ap-northeast-1",
		"eu-central-1",
		"us-west-1",
		"us-west-2",
	}

	path := "/var/tmp/hermes/reservation.out"
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return
	}

	repo := NewRepository(region)
	if err := repo.Fetch(); err != nil {
		t.Errorf("new repository: %v", err)
	}

	if err := repo.Write(path); err != nil {
		t.Errorf("write file: %v", err)
	}
}

func TestDeserialize(t *testing.T) {
	repo, err := Read("/var/tmp/hermes/reservation.out")
	if err != nil {
		t.Errorf("read file: %v", err)
	}

	if len(repo.SelectAll()) < 1 {
		t.Errorf("repository is empty")
	}
}

func TestRecommendM44xlarge(t *testing.T) {
	path := "/var/tmp/hermes/pricing/ap-northeast-1.out"
	repo, err := pricing.Read(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	rs := repo.FindByUsageType("APN1-BoxUsage:m4.4xlarge").
		OperatingSystem("Linux").
		Tenancy("Shared").
		PreInstalled("NA").
		LeaseContractLength("1yr").
		PurchaseOption("All Upfront").
		OfferingClass("standard")

	forecast := []pricing.Forecast{
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

	rec, _ := repo.Recommend(rs[0], forecast)
	min := rec.MinimumRecord

	if rs[0].OfferTermCode != min.OfferTermCode {
		t.Errorf("invalid offer term")
	}

	rsv, err := Read("/var/tmp/hermes/reservation.out")
	if err != nil {
		t.Errorf("read file: %v", err)
	}

	rs2 := rsv.FindByInstanceType(min.InstanceType).
		Region(min.Region).
		Duration(func(length string) int64 {
			duration := 31536000
			if length == "3yr" {
				duration = 94608000
			}
			return int64(duration)
		}(min.LeaseContractLength)).
		OfferingClass(min.OfferingClass).
		OfferingType(min.PurchaseOption).
		ProductDescription(min.OperatingSystem)

	if rs2[0].Count() != rs2[0].InstanceCount {
		t.Errorf("invalid count")
	}
}
