package pricing

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	region := c.String("region")
	dir := c.GlobalString("dir")
	format := c.String("format")

	repo, err := pricing.Read(fmt.Sprintf("%s/pricing/%s.out", dir, region))
	if err != nil {
		fmt.Printf("read pricing (region=%s): %v\n", region, err)
		os.Exit(1)
	}

	if format == "json" {
		for _, r := range repo.SelectAll() {
			fmt.Println(r.JSON())
		}
	}
}
