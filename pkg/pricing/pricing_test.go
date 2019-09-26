package pricing

import (
	"fmt"
	"sort"
	"testing"
)

func TestFetchRedshift(t *testing.T) {
	p, err := Fetch(Redshift, "ap-northeast-1")
	if err != nil {
		t.Error(err)
	}

	list := make([]Price, 0)
	for _, v := range p {
		list = append(list, v)
	}
	sort.SliceStable(list, func(i, j int) bool { return list[i].UsageType < list[j].UsageType })

	for i := range list {
		fmt.Println(list[i])
	}
}
