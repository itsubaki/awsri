package pricing

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	region := c.StringSlice("region")
	dir := c.GlobalString("dir")
	format := c.String("format")

	price, err := pricing.Deserialize(dir, region)
	if err != nil {
		fmt.Printf("deserialize: %v", err)
		os.Exit(1)
	}

	if format == "json" {
		bytes, err := json.Marshal(price)
		if err != nil {
			fmt.Printf("marshal: %v", err)
			os.Exit(1)
		}

		fmt.Println(string(bytes))
		return
	}

	if format == "csv" {
		for _, p := range price {
			fmt.Printf(
				"%s, %s, %s, %s, %s, %s, %s, %s, %s, %.3f, %.3f, %.3f\n",
				p.Version,
				p.Region,
				p.InstanceType,
				p.LeaseContractLength,
				p.PurchaseOption,
				p.OperatingSystem,
				p.CacheEngine,
				p.DatabaseEngine,
				p.OfferingClass,
				p.OnDemand,
				p.ReservedQuantity,
				p.ReservedHrs,
			)
		}
		return
	}
}
