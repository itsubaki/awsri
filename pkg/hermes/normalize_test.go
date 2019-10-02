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
	for k, _ := range family {
		if !strings.Contains(k, "BoxUsage:t2") {
			continue
		}
		if !strings.Contains(k, "Linux") {
			continue
		}
		if !strings.Contains(k, "NA") {
			continue
		}

		fmt.Println(k)
		//		fmt.Println(v)
	}

	mini := pricing.Minimum(family, plist)

	for k, v := range mini {
		if !strings.Contains(k, "BoxUsage:t2") {
			continue
		}
		if !strings.Contains(k, "Linux") {
			continue
		}

		fmt.Printf("%s %s %s\n", k, v.Price.UsageType, v.Minimum.UsageType)
	}

	forecast := []usage.Quantity{
		{
			Region:       "ap-northeast-1",
			UsageType:    "APN1-BoxUsage:t2.micro",
			Platform:     "Linux/UNIX",
			InstanceHour: 1,
			InstanceNum:  1,
		},
		{
			Region:       "ap-northeast-1",
			UsageType:    "APN1-BoxUsage:t2.small",
			Platform:     "Linux/UNIX",
			InstanceHour: 3,
			InstanceNum:  3,
		},
		{
			Region:       "ap-northeast-1",
			UsageType:    "APN1-BoxUsage:t2.nano",
			Platform:     "Linux/UNIX",
			InstanceHour: 5,
			InstanceNum:  5,
		},
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
			InstanceHour: 3,
			InstanceNum:  3,
		},
		{
			Region:       "ap-northeast-1",
			UsageType:    "APN1-BoxUsage:m4.4xlarge",
			Platform:     "Linux/UNIX",
			InstanceHour: 5,
			InstanceNum:  5,
		},
	}

	n := Normalize(forecast, mini)
	for _, nn := range n {
		fmt.Println(nn)
	}
}
