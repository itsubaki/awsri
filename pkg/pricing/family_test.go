package pricing

import (
	"fmt"
	"testing"
)

func TestFamily(t *testing.T) {
	plist, err := Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize pricing: %v", err)
	}

	family := Family(plist)
	for _, v := range family {
		fmt.Println(v)
	}
}
