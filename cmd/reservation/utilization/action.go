package utilization

import (
	"fmt"

	"github.com/itsubaki/hermes/calendar"
	"github.com/itsubaki/hermes/flag"
	"github.com/itsubaki/hermes/pricing"
	"github.com/itsubaki/hermes/reservation"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	dir := c.String("dir")
	format := c.String("format")
	normalize := c.Bool("normalize")
	merge := c.Bool("merge")
	groupby := c.Bool("groupby")
	period := c.String("period")
	attribute := c.String("attribute")
	region := flag.Split(c.StringSlice("region"))

	date, err := calendar.Last(period)
	if err != nil {
		return fmt.Errorf("get last period=%s: %v", period, err)
	}

	res, err := reservation.Deserialize(dir, date)
	if err != nil {
		return fmt.Errorf("deserialize: %v\n", err)
	}

	plist, err := pricing.Deserialize(dir, region)
	if err != nil {
		return fmt.Errorf("desirialize pricing: %v\n", err)
	}

	for _, e := range reservation.AddOnDemandConversionCost(plist, res) {
		fmt.Printf("[WARN] %s\n", e)
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

		return nil
	}

	if format == "json" && groupby {
		g, _ := reservation.GroupBy(res)
		for _, m := range g {
			fmt.Println(m)
		}

		return nil
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
					if attribute == "ondemand-conversion-cost" {
						fmt.Printf("%.3f, ", r.OnDemandConversionCost)
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

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}
