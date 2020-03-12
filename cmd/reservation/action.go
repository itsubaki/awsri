package reservation

import (
	"fmt"
	"math"
	"os"

	"github.com/itsubaki/hermes/pkg/usage"

	"github.com/itsubaki/hermes/pkg/pricing"

	"github.com/itsubaki/hermes/pkg/reservation"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	dir := c.GlobalString("dir")
	region := c.StringSlice("region")
	format := c.String("format")
	normalize := c.Bool("normalize")
	merge := c.Bool("merge")
	groupby := c.Bool("groupby")
	period := c.String("period")
	attribute := c.String("attribute")

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

	plist, err := pricing.Deserialize(dir, region)
	if err != nil {
		fmt.Printf("desirialize pricing: %v\n", err)
		os.Exit(1)
	}

	cache, warning := reservation.NewCache(plist)
	for _, w := range warning {
		fmt.Println(w)
	}

	for i := range res {
		p, err := cache.Find(res[i])
		if err != nil {
			fmt.Printf("[WARNING] %v\n", err)
			continue
		}

		res[i].CoveringCost = math.Round(p.OnDemand*res[i].Hours*1000) / 1000
	}

	if normalize {
		family := pricing.Family(plist)
		mini := pricing.Minimum(plist, family)
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
					if attribute == "covering-cost" {
						fmt.Printf("%.3f, ", r.CoveringCost)
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
