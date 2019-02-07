package linear

import (
	"github.com/itsubaki/hermes/cmd/hermes/recommend"
	"github.com/urfave/cli"
)

type Regression struct {
}

func (lr *Regression) Do(c *cli.Context) recommend.ForecastList {
	output := recommend.ForecastList{}
	return output
}
