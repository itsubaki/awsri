package pricing

import (
	"fmt"
	"strings"
	"testing"
)

func TestMinimum(t *testing.T) {
	plist, err := Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize pricing: %v", err)
	}

	family := Family(plist)
	mini := Minimum(family, plist)

	for k, v := range mini {
		if !strings.Contains(k, "BoxUsage:c4.8x") {
			continue
		}
		if !strings.Contains(k, "Linux") {
			continue
		}

		fmt.Printf("%s %s %s\n", k, v.Price.NormalizationSizeFactor, v.Minimum.NormalizationSizeFactor)
	}
}
