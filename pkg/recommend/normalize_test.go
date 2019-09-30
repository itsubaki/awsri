package recommend

import (
	"testing"

	"github.com/itsubaki/hermes/pkg/usage"

	"github.com/itsubaki/hermes/pkg/pricing"
)

func TestNormalize(t *testing.T) {
	price, err := pricing.Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize: %v", err)
	}

	quantity := usage.Quantity{
		Region:       "ap-northeast-1",
		UsageType:    "APN1-BoxUsage:c4.2xlarge",
		Platform:     "Linux/UNIX",
		InstanceHour: 518332.57223100006,
		InstanceNum:  719.9063503208333,
	}

	q, err := Normalize([]usage.Quantity{quantity}, price)
	if err != nil {
		t.Errorf("normalize: %v", err)
	}

	if q[0].UsageType != "APN1-BoxUsage:c4.large" {
		t.Errorf("%v", q)
	}

	if q[0].InstanceHour != quantity.InstanceHour*4 {
		t.Errorf("%v", q)
	}

	if q[0].InstanceNum != quantity.InstanceNum*4 {
		t.Errorf("%v", q)
	}
}
