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
	normalize := c.Bool("normalize")
	merge := c.Bool("merge")
	monthly := c.Bool("monthly")
	attribute := c.String("attribute")
	date := usage.LastNMonths(12)

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

	if merge {
		res = reservation.Merge(res)
	}

	if format == "json" && !monthly {
		reservation.Sort(res)
		for _, r := range res {
			fmt.Println(r)
		}
	}

	if format == "json" && monthly {
		m := reservation.Monthly(res)
		for _, mm := range m {
			fmt.Println(mm)
		}
		return
	}

	if format == "csv" {
		fmt.Printf("account_id, description, region, instance_type, usage_type, os/engine, deploymet_option, ")
		for i := range date {
			fmt.Printf("%s, ", date[i].YYYYMM())
		}
		fmt.Println()

		m := reservation.Monthly(res)
		keys := reservation.SortedKey(m)
		for _, k := range keys {
			fmt.Printf(
				"%s, %s, %s, %s, %s, %s, %s, ",
				m[k][0].AccountID,
				m[k][0].Description,
				m[k][0].Region,
				m[k][0].InstanceType,
				m[k][0].UsageType(),
				m[k][0].OSEngine(),
				m[k][0].DeploymentOption,
			)

			for _, d := range date {
				found := false
				for _, r := range m[k] {
					if d.YYYYMM() != r.Date {
						continue
					}

					if attribute == "num" {
						fmt.Printf("%.3f, ", r.Num)
					}
					if attribute == "hours" {
						fmt.Printf("%.3f, ", r.Hours)
					}
					if attribute == "percentage" {
						fmt.Printf("%.3f, ", r.Percentage)
					}

					found = true
					break
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
