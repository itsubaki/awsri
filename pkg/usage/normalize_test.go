package usage

import (
	"testing"

	"github.com/itsubaki/hermes/pkg/pricing"
)

func TestNormalize(t *testing.T) {
	plist, err := pricing.Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize pricing: %v", err)
	}

	family := pricing.Family(plist)
	mini := pricing.Minimum(family, plist)

	forecast := []Quantity{
		{
			Region:       "ap-northeast-1",
			UsageType:    "APN1-BoxUsage:m4.large",
			Platform:     "Linux/UNIX",
			InstanceHour: 1,
			InstanceNum:  1,
		},
		{
			Region:       "ap-northeast-1",
			UsageType:    "APN1-BoxUsage:m4.xlarge",
			Platform:     "Linux/UNIX",
			InstanceHour: 2,
			InstanceNum:  2,
		},
		{
			Region:       "ap-northeast-1",
			UsageType:    "APN1-BoxUsage:m4.2xlarge",
			Platform:     "Linux/UNIX",
			InstanceHour: 4,
			InstanceNum:  4,
		},
	}

	n := Normalize(forecast, mini)
	for _, nn := range n {
		if nn.UsageType != "APN1-BoxUsage:m4.large" {
			t.Errorf("%s\n", nn)
		}
	}
}
