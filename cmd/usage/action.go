package usage

import (
	"fmt"
	"sort"

	"github.com/itsubaki/hermes/pkg/calendar"
	"github.com/itsubaki/hermes/pkg/flag"
	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	dir := c.String("dir")
	format := c.String("format")
	normalize := c.Bool("normalize")
	merge := c.Bool("merge")
	overall := c.Bool("merge-overall")
	groupby := c.Bool("groupby")
	period := c.String("period")
	attribute := c.String("attribute")
	region := flag.Split(c.StringSlice("region"))

	date, err := calendar.Last(period)
	if err != nil {
		return fmt.Errorf("get last period=%s: %v", period, err)
	}

	quantity, err := usage.Deserialize(dir, date)
	if err != nil {
		return fmt.Errorf("deserialize usage: %v\n", err)
	}

	if normalize {
		plist, err := pricing.Deserialize(dir, region)
		if err != nil {
			return fmt.Errorf("desirialize pricing: %v\n", err)
		}

		family := pricing.Family(plist)
		mini := pricing.Minimum(plist, family)
		quantity = usage.Normalize(quantity, mini)
	}

	if merge {
		quantity = usage.Merge(quantity)
	}

	if overall {
		quantity = usage.MergeOverall(quantity)
	}

	if format == "json" && !groupby {
		usage.Sort(quantity)
		for _, q := range quantity {
			fmt.Println(q)
		}

		return nil
	}

	if format == "json" && groupby {
		g, _ := usage.GroupBy(quantity)
		for _, q := range g {
			fmt.Println(q)
		}

		return nil
	}

	if format == "csv" {
		sort.Slice(date, func(i, j int) bool { return date[i].Start < date[j].Start })

		fmt.Printf("account_id, description, region, usage_type, os/engine, ")
		for i := range date {
			fmt.Printf("%s, ", date[i].String())
		}
		fmt.Println()

		g, keys := usage.GroupBy(quantity)
		for _, k := range keys {
			fmt.Printf("%s, %s, ", g[k][0].AccountID, g[k][0].Description)
			fmt.Printf("%s, %s, %s, ", g[k][0].Region, g[k][0].UsageType, g[k][0].OSEngine())

			for _, d := range date {
				found := false
				for _, q := range g[k] {
					if d.YYYYMMDD() != q.Date {
						continue
					}

					if q.Unit == "Requests" {
						fmt.Printf("%d, ", q.Requests)
					}
					if q.Unit == "GB" {
						fmt.Printf("%.10f, ", q.GByte)
					}
					if q.Unit == "Hrs" && attribute == "num" {
						fmt.Printf("%.3f, ", q.InstanceNum)
					}
					if q.Unit == "Hrs" && attribute == "hours" {
						fmt.Printf("%.3f, ", q.InstanceHour)
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
