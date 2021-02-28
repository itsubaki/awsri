package recommend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/recommend"
	"github.com/itsubaki/hermes/pkg/usage"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	format := c.String("format")

	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("read stdin: %v", err)
	}

	type Purchase struct {
		Price    pricing.Price    `json:"price"`
		Quantity []usage.Quantity `json:"quantity"`
	}

	var purchase []Purchase
	if err := json.Unmarshal(stdin, &purchase); err != nil {
		return fmt.Errorf("unmarshal: %v\n", err)
	}

	out := make([]recommend.Recommended, 0)
	for _, p := range purchase {
		r := recommend.Do(p.Quantity, p.Price)
		out = append(out, r)
	}

	if format == "json" {
		for _, o := range out {
			fmt.Println(o)
		}

		return nil
	}

	if format == "csv" {
		fmt.Printf("usage_type, lease_contract_length, purchase_option, os/engine, ")
		fmt.Printf("total_hours, ondemand_hours, reserved_hours, reserved_instance_num, ")
		fmt.Printf("full_ondemand_cost, reserved_applied_cost, saving_cost")
		fmt.Println()

		for _, o := range out {
			fmt.Printf("%s, %s, %s, %s, %f, %f, %f, %d, %f, %f, %f\n",
				o.Price.UsageType,
				o.Price.LeaseContractLength,
				o.Price.PurchaseOption,
				o.Price.OSEngine(),
				o.Usage.TotalHours,
				o.Usage.OnDemandHours,
				o.Usage.ReservedHours,
				o.Usage.ReservedInstanceNum,
				o.Cost.FullOnDemand,
				o.Cost.ReservedApplied.Total,
				o.Cost.Saving,
			)
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}
