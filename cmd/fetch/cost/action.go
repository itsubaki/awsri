package cost

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/itsubaki/hermes/pkg/cost"
	"github.com/itsubaki/hermes/pkg/usage"

	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	dir := c.GlobalString("dir")

	path := fmt.Sprintf("%s/cost", dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	date := usage.LastNMonths(12)
	for i := range date {
		file := fmt.Sprintf("%s/%s.out", path, date[i].YYYYMM())
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			continue
		}

		ac, err := cost.Fetch(date[i].Start, date[i].End)
		if err != nil {
			fmt.Printf("fetch cost (%s, %s): %v\n", date[i].Start, date[i].End, err)
			os.Exit(1)
		}

		bytes, err := json.Marshal(ac)
		if err != nil {
			fmt.Printf("marshal: %v\n", err)
			os.Exit(1)
		}

		if err := ioutil.WriteFile(file, bytes, os.ModePerm); err != nil {
			fmt.Errorf("write file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("write: %v\n", file)
	}
}
