package predict

import (
	"encoding/json"
	"fmt"

	"github.com/itsubaki/hermes/cmd/hermes/recommend"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {

	output := recommend.ForecstList{}

	bytes, err := json.Marshal(&output)
	if err != nil {
		fmt.Println(fmt.Errorf("marshal: %v", err))
		return
	}

	fmt.Println(string(bytes))
}
