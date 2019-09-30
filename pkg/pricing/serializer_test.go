package pricing

import (
	"testing"
)

func TestSerialize(t *testing.T) {

}

func TestDeserialize(t *testing.T) {
	_, err := Deserialize("/var/tmp/hermes", []string{"ap-northeast-1"})
	if err != nil {
		t.Errorf("desirialize: %v", err)
	}
}
