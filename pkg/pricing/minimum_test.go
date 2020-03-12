package pricing

import (
	"fmt"
	"testing"
)

func TestMinimum(t *testing.T) {
	plist, err := Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize pricing: %v", err)
	}

	family := Family(plist)
	mini := Minimum(plist, family)
	for _, v := range mini {
		fmt.Println(v)
	}
}
