package fetch

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/pkg/calendar"

	"github.com/itsubaki/hermes/pkg/cost"
	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/reservation"
	"github.com/itsubaki/hermes/pkg/usage"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	dir := c.GlobalString("dir")
	region := c.StringSlice("region")
	period := c.String("period")

	date, err := calendar.Last(period)
	if err != nil {
		fmt.Printf("get last period=%s: %v", period, err)
		os.Exit(1)
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

	if err := pricing.Serialize(dir, region); err != nil {
		fmt.Printf("serialize pricing: %v", err)
		os.Exit(1)
	}
}
