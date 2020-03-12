package pricing

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	dir := c.GlobalString("dir")
	region := c.StringSlice("region")
	format := c.String("format")

	price, err := pricing.Deserialize(dir, region)
	if err != nil {
		fmt.Printf("deserialize: %v\n", err)
		os.Exit(1)
	}

	if format == "json" {
		for _, p := range price {
			fmt.Println(p)
		}
		return
	}

	if format == "csv" {
		fmt.Println("discount_rate, break_even_point(month), version, region, instance_type, usage_type, lease_contract_length, purchase_option, os/engine, tenancy, pre_installed, operation, offering_class, ondemand, reserved_quantity, reserved_hours, normalization_factor")
		for _, p := range price {
			fmt.Printf(
				"%.2f, %d, %s, %s, %s, %s, %s, %s, %s%s%s, %s, %s, %s, %s, %.3f, %.3f, %.3f, %s\n",
				p.DiscountRate(),
				p.BreakEvenPoint(),
				p.Version,
				p.Region,
				p.InstanceType,
				p.UsageType,
				p.LeaseContractLength,
				p.PurchaseOption,
				p.OperatingSystem,
				p.CacheEngine,
				p.DatabaseEngine,
				p.Tenancy,
				p.PreInstalled,
				p.Operation,
				p.OfferingClass,
				p.OnDemand,
				p.ReservedQuantity,
				p.ReservedHrs,
				p.NormalizationSizeFactor,
			)
		}
		return
	}
}
