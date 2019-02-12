package pricing

import "testing"

func TestUniqueOperatingSystem(t *testing.T) {
	path := "/var/tmp/hermes/pricing/ap-northeast-1.out"
	repo, err := Read(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	for _, r := range repo.SelectAll().Unique("OperatingSystem") {
		if r != "Windows" && r != "SUSE" && r != "RHEL" && r != "Linux" {
			t.Errorf("invalid OperatingSystem=%s", r)
		}
	}
}

func TestUniqueCacheEngine(t *testing.T) {
	path := "/var/tmp/hermes/pricing/ap-northeast-1.out"
	repo, err := Read(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	for _, r := range repo.SelectAll().Unique("CacheEngine") {
		if r != "Redis" && r != "Memcached" {
			t.Errorf("invalid CacheEngine=%s", r)
		}
	}
}

func TestUniqueDatabaseEngine(t *testing.T) {
	path := "/var/tmp/hermes/pricing/ap-northeast-1.out"
	repo, err := Read(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	for _, r := range repo.SelectAll().Unique("DatabaseEngine") {
		if r != "SQL Server" &&
			r != "PostgreSQL" &&
			r != "Oracle" &&
			r != "MySQL" &&
			r != "MariaDB" &&
			r != "Aurora PostgreSQL" &&
			r != "Aurora MySQL" {
			t.Errorf("invalid DatabaseEngine=%s", r)
		}
	}
}

func TestBreakevenPoint1yr(t *testing.T) {
	path := "/var/tmp/hermes/pricing/ap-northeast-1.out"
	repo, err := Read(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	rs := repo.FindByInstanceType("m4.large").
		OperatingSystem("Linux").
		Tenancy("Shared").
		PreInstalled("NA").
		OfferingClass("standard").
		LeaseContractLength("1yr").
		PurchaseOption("All Upfront")

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
	path := "/var/tmp/hermes/pricing/ap-northeast-1.out"
	repo, err := Read(path)
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
	path := "/var/tmp/hermes/pricing/ap-northeast-1.out"
	repo, err := Read(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	rs := repo.FindByInstanceType("cache.m4.large").
		CacheEngine("Redis").
		PurchaseOption("Heavy Utilization").
		LeaseContractLength("3yr")

	for _, r := range rs {
		if r.CacheEngine != "Redis" {
			t.Error("invalid engine")
		}

		e := r.GetCost(0, 10)
		if e.ReservedApplied.OnDemand != 0 {
			t.Error("invalid reserved applied")
		}

		if e.Difference < 0 {
			t.Error("invalid difference")
		}

		if e.DiscountRate < 0 {
			t.Error("invalid discount rate")
		}
	}
}

func TestRecommendNoReserved(t *testing.T) {
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
		{Date: "2018-06", InstanceNum: 10},
		{Date: "2018-07", InstanceNum: 20},
		{Date: "2018-08", InstanceNum: 10},
		{Date: "2018-09", InstanceNum: 20},
		{Date: "2018-10", InstanceNum: 10},
		{Date: "2018-11", InstanceNum: 20},
	}

	rec := r.Recommend(forecast)
	if rec.FullOnDemandCost != rec.ReservedAppliedCost.Total {
		t.Errorf("invalid total cost")
	}

	if rec.OnDemandInstanceNumAvg != 15 {
		t.Errorf("invalid ondemand instance num")
	}

	if rec.ReservedInstanceNum != 0 {
		t.Errorf("invalid reserved instance num")
	}
}

func TestRecommend1yr(t *testing.T) {
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

	forecast := []Forecast{
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

	rec := r.Recommend(forecast)
	if rec.OnDemandInstanceNumAvg != 23.7 {
		t.Errorf("invalid ondemand instance num")
	}

	if rec.ReservedInstanceNum != 50 {
		t.Errorf("invalid reserved instance num")
	}

	if rec.ReservedQuantity != 35650 {
		t.Errorf("invalid reserved quantity")
	}

	if rec.Difference != 20852.000000000007 {
		t.Error("invalid difference")
	}

	if rec.DiscountRate != 0.2503723766793573 {
		t.Error("invalid discount rate")
	}
}

func TestRecommend1yrMinimum(t *testing.T) {
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

	rec := r.Recommend(forecast, "minimum")
	if rec.ReservedInstanceNum != 10 {
		t.Errorf("invalid reserved instance num")
	}
}

func TestRecommend3yr(t *testing.T) {
	r := &Record{
		SKU:                     "7MYWT7Y96UT3NJ2D",
		OfferTermCode:           "NQ3QZPMQV9",
		Region:                  "ap-northeast-1",
		InstanceType:            "m4.large",
		UsageType:               "APN1-BoxUsage:m4.large",
		LeaseContractLength:     "3yr",
		PurchaseOption:          "All Upfront",
		OnDemand:                0.129,
		ReservedHrs:             0,
		ReservedQuantity:        1457,
		Tenancy:                 "Shared",
		PreInstalled:            "NA",
		OperatingSystem:         "Linux",
		Operation:               "RunInstances",
		OfferingClass:           "standard",
		NormalizationSizeFactor: "4",
	}

	forecast := []Forecast{
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
		{Date: "2019-01", InstanceNum: 120.4},
		{Date: "2019-02", InstanceNum: 110.3},
		{Date: "2019-03", InstanceNum: 100.1},
		{Date: "2019-04", InstanceNum: 90.9},
		{Date: "2019-05", InstanceNum: 80.9},
		{Date: "2019-06", InstanceNum: 70.6},
		{Date: "2019-07", InstanceNum: 60.3},
		{Date: "2019-08", InstanceNum: 50.9},
		{Date: "2019-09", InstanceNum: 40.7},
		{Date: "2019-10", InstanceNum: 30.6},
		{Date: "2019-11", InstanceNum: 20.2},
		{Date: "2019-12", InstanceNum: 10.8},
		{Date: "2020-01", InstanceNum: 120.4},
		{Date: "2020-02", InstanceNum: 110.3},
		{Date: "2020-03", InstanceNum: 100.1},
		{Date: "2020-04", InstanceNum: 90.9},
		{Date: "2020-05", InstanceNum: 80.9},
		{Date: "2020-06", InstanceNum: 70.6},
		{Date: "2020-07", InstanceNum: 60.3},
		{Date: "2020-08", InstanceNum: 50.9},
		{Date: "2020-09", InstanceNum: 40.7},
		{Date: "2020-10", InstanceNum: 30.6},
		{Date: "2020-11", InstanceNum: 20.2},
		{Date: "2020-12", InstanceNum: 10.8},
	}

	rec := r.Recommend(forecast)
	if rec.BreakevenPointInMonth != 16 {
		t.Errorf("failed 3yr recommend")
	}
}
