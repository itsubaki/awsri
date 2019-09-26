package recommend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/itsubaki/hermes/pkg/usage"
)

func TestRecommend(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	quantity := make([]usage.Quantity, 0)
	date := usage.Last12Months()
	for i := range date {
		file := fmt.Sprintf("/var/tmp/hermes/usage/%s.out", date[i].YYYYMM())
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("file not found: %v", file)
		}

		read, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("read %s: %v", file, err)
		}

		var q []usage.Quantity
		if err := json.Unmarshal(read, &q); err != nil {
			t.Errorf("unmarshal: %v", err)
		}

		quantity = append(quantity, q...)
	}

	res, err := Recommend(quantity)
	if err != nil {
		t.Errorf("recommend: %v", err)
	}

	for _, r := range res {
		fmt.Println(r)
	}
}
