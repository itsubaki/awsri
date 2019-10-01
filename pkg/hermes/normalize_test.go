package hermes

import (
	"testing"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func TestNormalize(t *testing.T) {
	plist, err := pricing.Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize: %v", err)
	}

	p := pricing.Price{
		Region:                  "ap-northeast-1",
		UsageType:               "APN1-BoxUsage:c4.2xlarge",
		Tenancy:                 "Shared",
		PreInstalled:            "NA",
		OperatingSystem:         "Linux",
		OfferingClass:           "standard",
		LeaseContractLength:     "1yr",
		PurchaseOption:          "All Upfront",
		OnDemand:                0.126 * 4,
		ReservedQuantity:        738 * 4,
		ReservedHrs:             0 * 4,
		NormalizationSizeFactor: "16",
	}

	quantity := usage.Quantity{
		Region:       "ap-northeast-1",
		UsageType:    "APN1-BoxUsage:c4.2xlarge",
		Platform:     "Linux/UNIX",
		InstanceHour: 518332.57223100006,
		InstanceNum:  719.9063503208333,
	}

	q, m, err := Normalize(quantity, p, plist)
	if err != nil {
		t.Errorf("normalize: %v", err)
	}

	if q.UsageType != "APN1-BoxUsage:c4.large" {
		t.Errorf("%v", q)
	}

	if q.InstanceHour != quantity.InstanceHour*4 {
		t.Errorf("%v", q)
	}

	if q.InstanceNum != quantity.InstanceNum*4 {
		t.Errorf("%v", q)
	}

	if m.UsageType != "APN1-BoxUsage:c4.large" {
		t.Errorf("%v", m)
	}
}
