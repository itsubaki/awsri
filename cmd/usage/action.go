package usage

import (
	"fmt"
	"os"
	"sort"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	region := c.StringSlice("region")
	dir := c.GlobalString("dir")
	format := c.String("format")
	normalize := c.Bool("normalize")
	merge := c.Bool("merge")
	overall := c.Bool("merge-overall")
	monthly := c.Bool("monthly")
	attribute := c.String("attribute")
	date := usage.LastNMonths(c.Int("months"))

	quantity, err := usage.Deserialize(dir, date)
	if err != nil {
		fmt.Printf("deserialize usage: %v\n", err)
		os.Exit(1)
	}

	if normalize {
		plist, err := pricing.Deserialize(dir, region)
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

	if format == "json" && !monthly {
		usage.Sort(quantity)
		for _, q := range quantity {
			fmt.Println(q)
		}
		return
	}

	if format == "json" && monthly {
		mq := usage.Monthly(quantity)
		for _, q := range mq {
			fmt.Println(q)
		}
		return
	}

	if format == "csv" {
		sort.Slice(date, func(i, j int) bool { return date[i].Start < date[j].Start })

		fmt.Printf("account_id, description, region, usage_type, os/engine, ")
		for i := range date {
			fmt.Printf("%s, ", date[i].YYYYMM())
		}
		fmt.Println()

		mq := usage.Monthly(quantity)
		keys := usage.SortedKey(mq)
		for _, k := range keys {
			fmt.Printf("%s, %s, ", mq[k][0].AccountID, mq[k][0].Description)
			fmt.Printf("%s, %s, %s, ", mq[k][0].Region, mq[k][0].UsageType, mq[k][0].OSEngine())

			for _, d := range date {
				found := false
				for _, q := range mq[k] {
					if d.YYYYMM() != q.Date {
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
