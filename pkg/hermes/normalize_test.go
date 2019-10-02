package hermes

import (
	"fmt"
	"strings"
	"testing"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func TestNormalize(t *testing.T) {
	plist, err := pricing.Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize pricing: %v", err)
	}

	family := pricing.Family(plist)
	mini := pricing.Minimum(family, plist)

	for k, v := range mini {
		if !strings.Contains(k, "BoxUsage:c4") {
			continue
		}
		if !strings.Contains(k, "Linux") {
			continue
		}

		fmt.Printf("%s %s %s\n", k, v.Price.NormalizationSizeFactor, v.Minimum.NormalizationSizeFactor)
	}

	forecast := []usage.Quantity{
		{
			Region:       "ap-northeast-1",
			UsageType:    "APN1-BoxUsage:c4.large",
			Platform:     "Linux/UNIX",
			InstanceHour: 1,
			InstanceNum:  1,
		},
		{
			Region:       "ap-northeast-1",
			UsageType:    "APN1-BoxUsage:c4.xlarge",
			Platform:     "Linux/UNIX",
			InstanceHour: 3,
			InstanceNum:  3,
		},
		{
			Region:       "ap-northeast-1",
			UsageType:    "APN1-BoxUsage:c4.2xlarge",
			Platform:     "Linux/UNIX",
			InstanceHour: 5,
			InstanceNum:  5,
		},
		{
			Region:       "ap-northeast-1",
			UsageType:    "APN1-BoxUsage:c4.4xlarge",
			Platform:     "Linux/UNIX",
			InstanceHour: 7,
			InstanceNum:  7,
		},
	}

	n := Normalize(forecast, mini)
	for _, nn := range n {
		fmt.Println(nn)
	}
}
