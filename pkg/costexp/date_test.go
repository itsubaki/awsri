package costexp

import "testing"

func TestGetCurrentDate(t *testing.T) {
	current := GetCurrentDate()

	if len(current) != 12 {
		t.Errorf("invalid current date=%v", current)
	}
}
