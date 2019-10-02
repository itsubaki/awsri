package pricing

import (
	"fmt"
	"strings"
	"testing"
)

func TestFamily(t *testing.T) {
	plist, err := Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize pricing: %v", err)
	}

	family := Family(plist)
	for k, v := range family {
		if !strings.Contains(k, "BoxUsage:c4") {
			continue
		}
		if !strings.Contains(k, "Linux") {
			continue
		}

		if v.NormalizationSizeFactor != "8" {
			continue
		}

		fmt.Printf("%s %s\n", k, v)
	}
}
