package reservation

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/pkg/reservation"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	dir := c.GlobalString("dir")
	format := c.String("format")
	monthly := c.Bool("monthly")

	date := reservation.Last12Months()
	res, err := reservation.Deserialize(dir, date)
	if err != nil {
		fmt.Printf("deserialize: %v\n", err)
		os.Exit(1)
	}

	if format == "json" && !monthly {
		for _, r := range res {
			fmt.Println(r)
		}
	}

	if format == "json" && monthly {
		mr := reservation.Monthly(res)
		for _, r := range mr {
			fmt.Println(r)
		}
		return
	}

	if format == "csv" {
		fmt.Printf("accountID, region, instance_type, os/engine, ")
		for i := range date {
			fmt.Printf("%s, ", date[i].YYYYMM())
		}
		fmt.Println()

		mr := reservation.Monthly(res)
		keys := reservation.SortedKey(mr)
		for _, k := range keys {
			fmt.Printf("%s, %s, %s, ", mr[k][0].AccountID, mr[k][0].Region, mr[k][0].InstanceType)
			fmt.Printf("%s, ", fmt.Sprintf("%s%s%s", mr[k][0].Platform, mr[k][0].CacheEngine, mr[k][0].DatabaseEngine))

			for _, d := range date {
				found := false
				for _, r := range mr[k] {
					if d.YYYYMM() == r.Date {
						fmt.Printf("%.3f, ", r.Hours)
						found = true
						break
					}
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
