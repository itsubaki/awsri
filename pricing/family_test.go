package pricing_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/pricing"
)

func TestFamily(t *testing.T) {
	plist, err := pricing.Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize pricing: %v", err)
	}

	family := pricing.Family(plist)
	for _, v := range family {
		fmt.Println(v)
	}
}
