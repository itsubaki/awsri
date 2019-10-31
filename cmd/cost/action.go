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
	attribute := c.String("attribute")
	date := usage.LastNMonths(c.Int("months"))

	ac, err := cost.Deserialize(dir, date)
	if err != nil {
		fmt.Errorf("deserialize cost: %v\n", err)
		os.Exit(1)
	}

	if format == "json" {
		for _, a := range ac {
			fmt.Println(a)
		}
	}

	if format == "csv" {
		fmt.Printf("account_id, description, ")
		for i := range date {
			fmt.Printf("%s, ", date[i].YYYYMM())
		}
		fmt.Println()

		mc := cost.Monthly(ac)
		keys := cost.SortedKey(mc)
		for _, k := range keys {
			fmt.Printf("%s, %s, ", mc[k][0].AccountID, mc[k][0].Description)

			for _, d := range date {
				found := false
				for _, q := range mc[k] {
					if d.YYYYMM() != q.Date {
						continue
					}

					if attribute == "amortized" {
						fmt.Printf("%s, ", q.AmortizedCost.Amount)
					}
					if attribute == "net-amortized" {
						fmt.Printf("%s, ", q.NetAmortizedCost.Amount)
					}
					if attribute == "blended" {
						fmt.Printf("%s, ", q.BlendedCost.Amount)
					}
					if attribute == "unblended" {
						fmt.Printf("%s, ", q.UnblendedCost.Amount)
					}
					if attribute == "net-unblended" {
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
		return
	}
}
