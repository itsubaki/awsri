package predict

import (
	"encoding/json"
	"fmt"

	"github.com/itsubaki/hermes/cmd/hermes/predict/linear"
	"github.com/itsubaki/hermes/cmd/hermes/recommend"
	"github.com/urfave/cli"
)

type Predict interface {
	Do(c *cli.Context) recommend.ForecastList
}

func Action(c *cli.Context) {
	predict := &linear.Regression{}

	output := predict.Do(c)
	bytes, err := json.Marshal(output)
	if err != nil {
		fmt.Println(fmt.Errorf("marshal: %v", err))
		return
	}

	fmt.Println(string(bytes))
}
