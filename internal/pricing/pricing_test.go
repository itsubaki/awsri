package pricing

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestReadEC2Price(t *testing.T) {
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

func TestReadCachePrice(t *testing.T) {
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

func TestReadRDSPrice(t *testing.T) {
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

func TestGetEC2Price(t *testing.T) {
	p, err := Fetch(EC2URL, "ap-northeast-1")
	if err != nil {
		t.Error(err)
	}

	for k, v := range p {
		fmt.Printf("%v -> %v\n", k, v)
	}
}

func TestGetCachePrice(t *testing.T) {
	p, err := Fetch(CacheURL, "ap-northeast-1")
	if err != nil {
		t.Error(err)
	}

	for k, v := range p {
		fmt.Printf("%v -> %v\n", k, v)
	}
}

func TestGetRDSPrice(t *testing.T) {
	p, err := Fetch(RDSURL, "ap-northeast-1")
	if err != nil {
		t.Error(err)
	}

	for k, v := range p {
		fmt.Printf("%v -> %v\n", k, v)
	}
}
