package reserved

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/pkg/reservation/reserved"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	region := c.StringSlice("region")
	format := c.String("format")

	rsv, err := reserved.Fetch(region)
	if err != nil {
		fmt.Printf("fetch region=%s: %v", region, err)
		os.Exit(1)
	}

	if format == "json" {
		for _, r := range rsv {
			fmt.Println(r)
		}
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
	}
}
