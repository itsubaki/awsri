package recommend

import (
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/pkg/usage"
)

func TestMerge(t *testing.T) {
	quantity := []usage.Quantity{
		{
			UsageType:    "APN1-BoxUsage:c4.large",
			InstanceNum:  100,
			InstanceHour: 1000,
		},
		{
			UsageType:    "APN1-BoxUsage:c4.large",
			InstanceNum:  200,
			InstanceHour: 2000,
		},
		{
			UsageType:    "APN1-BoxUsage:c4.2xlarge",
			InstanceNum:  100,
			InstanceHour: 3000,
		},
	}

	q := Merge(quantity)
	if len(q) != 2 {
		t.Errorf("len: %v", q)
	}

	for _, qq := range q {
		fmt.Printf("%#v\n", qq)
	}
}
