package hermes

import (
	"fmt"
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

	//fmt.Println("family map------")
	//for k := range family {
	//	if !strings.Contains(k, "BoxUsage:m4") {
	//		continue
	//	}
	//	if !strings.Contains(k, "Linux") {
	//		continue
	//	}
	//	if !strings.Contains(k, "NA") {
	//		continue
	//	}
	//
	//	fmt.Printf("%s -> %s\n", k, family[k])
	//}
	//
	mini := pricing.Minimum(family, plist)

	//fmt.Println("mini map------")
	//for k := range mini {
	//	if !strings.Contains(k, "BoxUsage:m4.2x") {
	//		continue
	//	}
	//	if !strings.Contains(k, "Linux") {
	//		continue
	//	}
	//
	//	fmt.Printf("%s -> %s\n", mini[k].Price, mini[k].Minimum)
	//}
	//
	forecast := []usage.Quantity{
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

	fmt.Println("------")
	for _, nn := range n {
		fmt.Println(nn)
	}
}
