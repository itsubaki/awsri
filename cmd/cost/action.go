package cost

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/pkg/cost"
	"github.com/itsubaki/hermes/pkg/usage"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	dir := c.GlobalString("dir")
	format := c.String("format")
	date := usage.LastNMonths(c.Int("months"))

	ac, err := cost.Deserialize(dir, date)
	if err != nil {
		fmt.Printf("deserialize cost: %v\n", err)
		os.Exit(1)
	}

	if format == "json" {
		for _, a := range ac {
			fmt.Println(a)
		}
	}

	if format == "csv" {
		fmt.Printf("account_id, description, metric, ")
		for i := range date {
			fmt.Printf("%s, ", date[i].YYYYMM())
		}
		fmt.Println()

		mc := cost.Monthly(ac)
		keys := cost.SortedKey(mc)
		for _, k := range keys {
			for _, m := range []string{"amortized", "net-amortized", "blended", "unblended", "net-unblended"} {
				fmt.Printf("%s, %s, %s, ", mc[k][0].AccountID, mc[k][0].Description, m)

				for _, d := range date {
					found := false
					for _, q := range mc[k] {
						if d.YYYYMM() != q.Date {
							continue
						}

						if m == "amortized" {
							fmt.Printf("%s, ", q.AmortizedCost.Amount)
						}
						if m == "net-amortized" {
							fmt.Printf("%s, ", q.NetAmortizedCost.Amount)
						}
						if m == "blended" {
							fmt.Printf("%s, ", q.BlendedCost.Amount)
						}
						if m == "unblended" {
							fmt.Printf("%s, ", q.UnblendedCost.Amount)
						}
						if m == "net-unblended" {
							fmt.Printf("%s, ", q.NetUnblendedCost.Amount)
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
		}
		return
	}
}
