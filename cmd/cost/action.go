package cost

import (
	"fmt"

	"github.com/itsubaki/hermes/pkg/calendar"
	"github.com/itsubaki/hermes/pkg/cost"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	dir := c.String("dir")
	format := c.String("format")
	period := c.String("period")
	attribute := c.String("attribute")

	date, err := calendar.Last(period)
	if err != nil {
		return fmt.Errorf("get last period=%s: %v", period, err)
	}

	ac, err := cost.Deserialize(dir, date)
	if err != nil {
		return fmt.Errorf("deserialize cost: %v\n", err)
	}

	if format == "json" {
		for _, a := range ac {
			fmt.Println(a)
		}

		return nil
	}

	if format == "csv" {
		fmt.Printf("account_id, description, service, record_type, ")
		for i := range date {
			fmt.Printf("%s, ", date[i].String())
		}
		fmt.Println()

		g, keys := cost.GroupBy(ac)
		for _, k := range keys {
			fmt.Printf("%s, %s, %s, %s, ", g[k][0].AccountID, g[k][0].Description, g[k][0].Service, g[k][0].RecordType)

			for _, d := range date {
				found := false
				for _, q := range g[k] {
					if d.YYYYMMDD() != q.Date {
						continue
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
					if attribute == "amortized" {
						fmt.Printf("%s, ", q.AmortizedCost.Amount)
					}
					if attribute == "net-amortized" {
						fmt.Printf("%s, ", q.NetAmortizedCost.Amount)
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
