package recommend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/recommend"
	"github.com/itsubaki/hermes/pkg/usage"

	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("read stdin: %v", err)
		os.Exit(1)
	}

	type Purchase struct {
		Price    pricing.Price    `json:"price"`
		Quantity []usage.Quantity `json:"quantity"`
	}

	var purchase []Purchase
	if err := json.Unmarshal(stdin, &purchase); err != nil {
		fmt.Printf("unmarshal: %v\n", err)
		os.Exit(1)
	}

	for _, p := range purchase {
		r := recommend.Do(p.Quantity, p.Price)
		fmt.Println(r)
	}
}
