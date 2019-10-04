package fetch

import (
	"github.com/itsubaki/hermes/cmd/fetch/pricing"
	"github.com/itsubaki/hermes/cmd/fetch/usage"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	pricing.Action(c)
	usage.Action(c)
}
