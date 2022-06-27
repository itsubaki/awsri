package reserved

import (
	"fmt"

	"github.com/itsubaki/hermes/flag"
	"github.com/itsubaki/hermes/reservation/reserved"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	format := c.String("format")
	region := flag.Split(c.StringSlice("region"))

	rsv, err := reserved.Fetch(region)
	if err != nil {
		return fmt.Errorf("fetch region=%s: %v", region, err)
	}

	if format == "json" {
		for _, r := range rsv {
			fmt.Println(r)
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("reserved_id, region, duration, offering_type, offering_class, product_description, instance_type/class, count, multi-az, start, state")
		for _, r := range rsv {
			fmt.Printf(
				"%s, %s, %d, %s, %s, %s, %s, %d, %t, %s, %s",
				r.ReservedID,
				r.Region,
				r.Duration,
				r.OfferingType,
				r.OfferingClass,
				r.ProductDescription,
				r.TypeClass(),
				r.Count(),
				r.MultiAZ,
				r.Start,
				r.State,
			)
			fmt.Println()
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}
