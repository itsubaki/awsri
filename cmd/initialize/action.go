package initialize

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/pkg/billing"
	"github.com/itsubaki/hermes/pkg/costexp"
	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/reserved"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	region := c.StringSlice("region")
	dir := c.GlobalString("dir")

	if err := pricing.Download(region, dir); err != nil {
		fmt.Printf("write pricing: %v", err)
		os.Exit(1)
	}

	if err := costexp.Download(dir); err != nil {
		fmt.Printf("write costexp: %v", err)
		os.Exit(1)
	}

	if err := reserved.Download(region, dir); err != nil {
		fmt.Printf("write reservation: %v", err)
		os.Exit(1)
	}

	if err := billing.Download(dir); err != nil {
		fmt.Printf("write billing: %v", err)
		os.Exit(1)
	}
}
