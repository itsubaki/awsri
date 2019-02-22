package pricing

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"testing"
)

func TestReadCompute(t *testing.T) {
	path := fmt.Sprintf("%s/%s", os.Getenv("GOPATH"),
		"src/github.com/itsubaki/hermes/internal/pricing/_json/ec2/ap-northeast-1.json",
	)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("read body: %v", err)
	}

	price, err := Read("ap-northeast-1", buf)
	if err != nil {
		t.Error(err)
	}

	for _, v := range price {
		fmt.Printf("%v\n", v)
	}
}

func TestReadCache(t *testing.T) {
	path := fmt.Sprintf("%s/%s", os.Getenv("GOPATH"),
		"src/github.com/itsubaki/hermes/internal/pricing/_json/cache/ap-northeast-1.json",
	)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("read body: %v", err)
	}

	price, err := Read("ap-northeast-1", buf)
	if err != nil {
		t.Error(err)
	}

	for _, v := range price {
		fmt.Printf("%v\n", v)
	}
}

func TestReadDatabase(t *testing.T) {
	path := fmt.Sprintf("%s/%s", os.Getenv("GOPATH"),
		"src/github.com/itsubaki/hermes/internal/pricing/_json/rds/ap-northeast-1.json",
	)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("read body: %v", err)
	}

	price, err := Read("ap-northeast-1", buf)
	if err != nil {
		t.Error(err)
	}

	for _, v := range price {
		fmt.Printf("%v\n", v)
	}
}

func TestGetCompute(t *testing.T) {
	p, err := Fetch(ComputeURL, "ap-northeast-1")
	if err != nil {
		t.Error(err)
	}

	for k, v := range p {
		fmt.Printf("%v -> %v\n", k, v)
	}
}

func TestGetCache(t *testing.T) {
	p, err := Fetch(CacheURL, "ap-northeast-1")
	if err != nil {
		t.Error(err)
	}

	for k, v := range p {
		fmt.Printf("%v -> %v\n", k, v)
	}
}

func TestGetDatabase(t *testing.T) {
	p, err := Fetch(DatabaseURL, "ap-northeast-1")
	if err != nil {
		t.Error(err)
	}

	for k, v := range p {
		fmt.Printf("%v -> %v\n", k, v)
	}
}

func TestGetRedshift(t *testing.T) {
	p, err := Fetch(RedshiftURL, "ap-northeast-1")
	if err != nil {
		t.Error(err)
	}

	list := []OutputPrice{}
	for _, v := range p {
		list = append(list, v)
	}
	sort.SliceStable(list, func(i, j int) bool { return list[i].UsageType < list[j].UsageType })

	for i := range list {
		fmt.Println(list[i])
	}
}
