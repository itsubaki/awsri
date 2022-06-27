package pricing_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/hermes/pricing"
)

func TestSerialize(t *testing.T) {

}

func TestDeserialize(t *testing.T) {
	price, err := pricing.Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize: %v", err)
	}

	tmp := make(map[string]bool)
	for _, p := range price {
		if len(p.OperatingSystem) > 0 {
			tmp[p.OperatingSystem] = true
		}

		if len(p.CacheEngine) > 0 {
			tmp[p.CacheEngine] = true

		}

		if len(p.DatabaseEngine) > 0 {
			tmp[p.DatabaseEngine] = true
		}
	}

	for k := range tmp {
		fmt.Println(k)
	}
}
