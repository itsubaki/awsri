package fetch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	region := c.StringSlice("region")
	dir := c.GlobalString("dir")

	for _, r := range region {
		price := make([]pricing.Price, 0)
		for _, url := range pricing.URL {
			p, err := fetch(url, r)
			if err != nil {
				os.Exit(1)
			}
			price = append(price, p...)
		}

		path := fmt.Sprintf("%s/pricing", dir)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
		}

		bytes, err := json.Marshal(price)
		if err != nil {
			fmt.Printf("marshal: %v", err)
			os.Exit(1)
		}

		file := fmt.Sprintf("%s/%s.out", path, r)
		if err := ioutil.WriteFile(file, bytes, os.ModePerm); err != nil {
			fmt.Errorf("write file: %v", err)
			os.Exit(1)
		}

		fmt.Printf("write: %v\n", file)
	}
}

func fetch(url string, region string) ([]pricing.Price, error) {
	price, err := pricing.Fetch(url, region)
	if err != nil {
		return []pricing.Price{}, fmt.Errorf("fetch pricing (%s, %s): %v", url, region, err)
	}

	out := make([]pricing.Price, 0)
	for k := range price {
		out = append(out, price[k])
	}

	return out, nil
}
