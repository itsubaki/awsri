package fetch

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/itsubaki/hermes/pkg/cost"
	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/reservation"
	"github.com/itsubaki/hermes/pkg/usage"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	dir := c.GlobalString("dir")
	region := c.String("region")
	period := c.String("period")

	n, err := strconv.Atoi(period[:len(period)-1])
	if err != nil {
		fmt.Printf("invalid period(%v): %v", period, err)
		os.Exit(1)
	}

	var date []usage.Date
	if strings.HasSuffix(period, "m") {
		date = usage.LastNMonths(n)
	}
	if strings.HasSuffix(period, "d") {
		date = usage.LastNDays(n)
	}

	if err := cost.Serialize(dir, date); err != nil {
		fmt.Printf("serialize cost: %v", err)
		os.Exit(1)
	}

	if err := reservation.Serialize(dir, date); err != nil {
		fmt.Printf("serialize reservation: %v", err)
		os.Exit(1)
	}

	if err := usage.Serialize(dir, date); err != nil {
		fmt.Printf("serialize usage: %v", err)
		os.Exit(1)
	}

	if err := pricing.Serialize(dir, []string{region}); err != nil {
		fmt.Printf("serialize pricing: %v", err)
		os.Exit(1)
	}
}
