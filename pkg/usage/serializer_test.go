package usage

import (
	"fmt"
	"testing"
)

func TestDeserialize(t *testing.T) {
	usage, err := Deserialize("/var/tmp/hermes", LastNMonths(3))
	if err != nil {
		t.Errorf("desirialize: %v", err)
	}

	tmp := make(map[string]bool)
	for _, u := range usage {
		if len(u.Platform) > 0 {
			tmp[u.Platform] = true
		}

		if len(u.CacheEngine) > 0 {
			tmp[u.CacheEngine] = true

		}

		if len(u.DatabaseEngine) > 0 {
			tmp[u.DatabaseEngine] = true
		}
	}

	for k := range tmp {
		fmt.Println(k)
	}
}
