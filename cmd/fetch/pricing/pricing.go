package pricing

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	region := c.StringSlice("region")
	dir := c.GlobalString("dir")

	path := fmt.Sprintf("%s/pricing", dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	for _, r := range region {
		file := fmt.Sprintf("%s/%s.out", path, r)
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			continue
		}

		price := make([]pricing.Price, 0)
		for _, url := range pricing.URL {
			p, err := pricing.Fetch(url, r)
			if err != nil {
				fmt.Printf("fetch pricing (%s, %s): %v", url, r, err)
				os.Exit(1)
			}

			list := make([]pricing.Price, 0)
			for k := range p {
				list = append(list, p[k])
			}

			price = append(price, list...)
		}

		if err := pricing.Serialize(dir, r, price); err != nil {
			fmt.Printf("serialize: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("write: %v\n", file)
	}
}
