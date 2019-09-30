package recommend

import (
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/pkg/usage"
)

func TestMerge(t *testing.T) {
	quantity := []usage.Quantity{
		{
			UsageType:    "foobar",
			InstanceNum:  100,
			InstanceHour: 1000,
		},
		{
			UsageType:    "foobar",
			InstanceNum:  200,
			InstanceHour: 2000,
		},
	}

	q := Merge(quantity)
	for _, qq := range q {
		fmt.Printf("%#v\n", qq)
	}
}
