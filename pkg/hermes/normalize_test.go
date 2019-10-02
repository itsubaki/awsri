package hermes

import (
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func TestNormalize(t *testing.T) {
	// price list
	plist, err := pricing.Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize: %v", err)
	}

	family := pricing.Family(plist)
	mini := pricing.Minimum(family, plist)

	forecast := []usage.Quantity{
		{
			Region:       "ap-northeast-1",
			UsageType:    "APN1-BoxUsage:c4.2xlarge",
			Platform:     "Linux/UNIX",
			InstanceHour: 518332.57223100006,
			InstanceNum:  719.9063503208333,
		},
	}

	n := Normalize(forecast, mini)
	for _, nn := range n {
		fmt.Println(nn)
	}

	if n[0].InstanceHour != forecast[0].InstanceHour*4 {
		t.Errorf("%v", n[0])
	}

	if n[0].InstanceNum != forecast[0].InstanceNum*4 {
		t.Errorf("%v", n[0])
	}

	if n[0].UsageType != "APN1-BoxUsage:c4.large" {
		t.Errorf("%v", n[0])
	}
}
