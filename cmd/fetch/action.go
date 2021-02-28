package fetch

import (
	"fmt"

	"github.com/itsubaki/hermes/pkg/calendar"
	"github.com/itsubaki/hermes/pkg/cost"
	"github.com/itsubaki/hermes/pkg/flag"
	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/reservation"
	"github.com/itsubaki/hermes/pkg/usage"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	dir := c.String("dir")
	period := c.String("period")
	region := flag.Split(c.StringSlice("region"))
	metrics := flag.Split(c.StringSlice("metrics"))

	date, err := calendar.Last(period)
	if err != nil {
		return fmt.Errorf("get last period=%s: %v", period, err)
	}

	if err := cost.Serialize(dir, date, metrics); err != nil {
		return fmt.Errorf("serialize cost: %v", err)
	}

	if err := reservation.Serialize(dir, date); err != nil {
		return fmt.Errorf("serialize reservation: %v", err)
	}

	if err := usage.Serialize(dir, date); err != nil {
		return fmt.Errorf("serialize usage: %v", err)
	}

	if err := pricing.Serialize(dir, region); err != nil {
		return fmt.Errorf("serialize pricing: %v", err)
	}

	return nil
}
