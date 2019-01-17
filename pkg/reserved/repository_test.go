package reserved

import (
	"fmt"
	"strings"
	"testing"

	"github.com/itsubaki/hermes/pkg/awsprice"
)

func TestRepository(t *testing.T) {
	dir := "/var/tmp/hermes/reserved"
	path := fmt.Sprintf("%s/%s", dir, "example.out")
	repo, err := NewRepository(path)
	if err != nil {
		t.Errorf("new repository: %v", err)
	}

	if len(repo.SelectAll()) < 1 {
		t.Errorf("repository is empty")
	}

	if repo.Profile != "example" {
		t.Errorf("invalid profile")
	}
}

func TestGetReserved(t *testing.T) {
	repo, err := NewRepository("/var/tmp/hermes/reserved/example.out")
	if err != nil {
		t.Errorf("new repository: %v", err)
	}

	r := repo.SelectAll()[0]
	yr := "1yr"
	if r.Duration == 94608000 {
		yr = "3yr"
	}

	os := "Linux"
	if strings.Contains(r.ProductDescription, "Windows") {
		os = "Windows"
	}

	{
		repo, err := awsprice.NewRepository("/var/tmp/hermes/awsprice/ap-northeast-1.out")
		if err != nil {
			t.Errorf("%v", err)
		}

		rs := repo.FindByInstanceType(r.InstanceType).
			OfferingClass(r.OfferingClass).
			PurchaseOption(r.OfferingType).
			OperatingSystem(os).
			LeaseContractLength(yr)

		if len(rs) != 1 {
			t.Errorf("invalid resultset length")
		}
	}
}
