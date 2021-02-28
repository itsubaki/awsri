package pricing

import (
	"fmt"

	"github.com/itsubaki/hermes/pkg/flag"
	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	dir := c.String("dir")
	format := c.String("format")
	region := flag.Split(c.StringSlice("region"))

	price, err := pricing.Deserialize(dir, region)
	if err != nil {
		return fmt.Errorf("deserialize: %v\n", err)
	}

	if format == "json" {
		for _, p := range price {
			fmt.Println(p)
		}

		return nil
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

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}
