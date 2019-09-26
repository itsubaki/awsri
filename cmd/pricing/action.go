package pricing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	region := c.StringSlice("region")
	dir := c.GlobalString("dir")
	format := c.String("format")

	price := make([]pricing.Price, 0)
	for _, r := range region {
		file := fmt.Sprintf("%s/pricing/%s.out", dir, r)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("file not found: %v", file)
			os.Exit(1)
		}

		read, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("read %s: %v", file, err)
			os.Exit(1)
		}

		var p []pricing.Price
		if err := json.Unmarshal(read, &p); err != nil {
			fmt.Printf("unmarshal: %v", err)
			os.Exit(1)
		}

		price = append(price, p...)
	}

	sort.SliceStable(price, func(i, j int) bool { return price[i].Version < price[j].Version })
	sort.SliceStable(price, func(i, j int) bool { return price[i].Region < price[j].Region })
	sort.SliceStable(price, func(i, j int) bool { return price[i].InstanceType < price[j].InstanceType })
	sort.SliceStable(price, func(i, j int) bool { return price[i].LeaseContractLength < price[j].LeaseContractLength })
	sort.SliceStable(price, func(i, j int) bool { return price[i].PurchaseOption < price[j].PurchaseOption })

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
