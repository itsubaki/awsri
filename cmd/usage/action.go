package usage

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	region := c.String("region")
	dir := c.GlobalString("dir")
	format := c.String("format")
	normalize := c.Bool("normalize")
	merge := c.Bool("merge")
	overall := c.Bool("merge-overall")
	groupby := c.Bool("groupby")
	attribute := c.String("attribute")
	period := c.String("period")

	n, err := strconv.Atoi(period[:len(period)-1])
	if err != nil {
		fmt.Printf("invalid period(%v): %v", period, err)
		os.Exit(1)
	}

	var date []usage.Date
	if strings.HasSuffix(period, "m") {
		date = usage.LastNMonths(n)
	}
	if strings.HasSuffix(period, "d") {
		date = usage.LastNDays(n)
	}

	quantity, err := usage.Deserialize(dir, date)
	if err != nil {
		fmt.Printf("deserialize usage: %v\n", err)
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
		return
	}

	if format == "json" && groupby {
		g, _ := usage.GroupBy(quantity)
		for _, q := range g {
			fmt.Println(q)
		}
		return
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
		return
	}
}
