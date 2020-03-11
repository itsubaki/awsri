package reservation

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/pkg/usage"

	"github.com/itsubaki/hermes/pkg/pricing"

	"github.com/itsubaki/hermes/pkg/reservation"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	region := c.String("region")
	dir := c.GlobalString("dir")
	format := c.String("format")
	normalize := c.Bool("normalize")
	merge := c.Bool("merge")
	groupby := c.Bool("groupby")
	attribute := c.String("attribute")
	period := c.String("period")

	date, err := usage.Last(period)
	if err != nil {
		fmt.Printf("get last months/days: %v", err)
		os.Exit(1)
	}

	res, err := reservation.Deserialize(dir, date)
	if err != nil {
		fmt.Printf("deserialize: %v\n", err)
		os.Exit(1)
	}

	if normalize {
		plist, err := pricing.Deserialize(dir, []string{region})
		if err != nil {
			fmt.Printf("desirialize pricing: %v\n", err)
			os.Exit(1)
		}

		family := pricing.Family(plist)
		mini := pricing.Minimum(family, plist)

		res = reservation.Normalize(res, mini)
	}

	if merge {
		res = reservation.Merge(res)
	}

	if format == "json" && !groupby {
		reservation.Sort(res)
		for _, r := range res {
			fmt.Println(r)
		}
	}

	if format == "json" && groupby {
		g, _ := reservation.GroupBy(res)
		for _, m := range g {
			fmt.Println(m)
		}
		return
	}

	if format == "csv" {
		fmt.Printf("account_id, description, region, instance_type, usage_type, os/engine, deploymet_option, ")
		for i := range date {
			fmt.Printf("%s, ", date[i].String())
		}
		fmt.Println()

		g, keys := reservation.GroupBy(res)
		for _, k := range keys {
			fmt.Printf(
				"%s, %s, %s, %s, %s, %s, %s, ",
				g[k][0].AccountID,
				g[k][0].Description,
				g[k][0].Region,
				g[k][0].InstanceType,
				g[k][0].UsageType(),
				g[k][0].OSEngine(),
				g[k][0].DeploymentOption,
			)

			for _, d := range date {
				found := false
				for _, r := range g[k] {
					if d.YYYYMMDD() != r.Date {
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
