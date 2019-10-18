package reservation

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/pkg/pricing"

	"github.com/itsubaki/hermes/pkg/reservation"
	"github.com/itsubaki/hermes/pkg/usage"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	region := c.StringSlice("region")
	dir := c.GlobalString("dir")
	format := c.String("format")
	monthly := c.Bool("monthly")
	normalize := c.Bool("normalize")

	date := usage.Last12Months()
	res, err := reservation.Deserialize(dir, date)
	if err != nil {
		fmt.Printf("deserialize: %v\n", err)
		os.Exit(1)
	}

	if normalize {
		plist, err := pricing.Deserialize(dir, region)
		if err != nil {
			fmt.Errorf("desirialize pricing: %v\n", err)
			os.Exit(1)
		}

		family := pricing.Family(plist)
		mini := pricing.Minimum(family, plist)

		res = reservation.Normalize(res, mini)
	}

	if format == "json" && !monthly {
		for _, r := range res {
			fmt.Println(r)
		}
	}

	if format == "json" && monthly {
		mr := reservation.Monthly(res)
		for _, r := range mr {
			fmt.Println(r)
		}
		return
	}

	if format == "csv" {
		fmt.Printf("account_id, description, region, instance_type, os/engine, deploymet_option, ")
		for i := range date {
			fmt.Printf("%s, ", date[i].YYYYMM())
		}
		fmt.Println()

		mr := reservation.Monthly(res)
		keys := reservation.SortedKey(mr)
		for _, k := range keys {
			fmt.Printf(
				"%s, %s, %s, %s, %s, %s, ",
				mr[k][0].AccountID,
				mr[k][0].Description,
				mr[k][0].Region,
				mr[k][0].InstanceType,
				mr[k][0].OSEngine(),
				mr[k][0].DeploymentOption,
			)

			for _, d := range date {
				found := false
				for _, r := range mr[k] {
					if d.YYYYMM() == r.Date {
						fmt.Printf("%.3f, ", r.Hours)
						found = true
						break
					}
				}

				if !found {
					fmt.Printf("0.0, ")
				}
			}

			fmt.Println()
		}
		return
	}
}
