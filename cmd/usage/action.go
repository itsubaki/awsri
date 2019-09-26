package usage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/itsubaki/hermes/pkg/usage"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	dir := c.GlobalString("dir")
	format := c.String("format")

	quantity := make([]usage.Quantity, 0)
	date := usage.Last12Months()
	for i := range date {
		file := fmt.Sprintf("%s/usage/%s.out", dir, date[i].YYYYMM())
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("file not found: %v", file)
			os.Exit(1)
		}

		read, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("read %s: %v", file, err)
			os.Exit(1)
		}

		var q []usage.Quantity
		if err := json.Unmarshal(read, &q); err != nil {
			fmt.Printf("unmarshal: %v", err)
			os.Exit(1)
		}

		quantity = append(quantity, q...)
	}

	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].AccountID < quantity[j].AccountID })
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].Region < quantity[j].Region })
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].UsageType < quantity[j].UsageType })
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].Date < quantity[j].Date })

	if format == "json" {
		bytes, err := json.Marshal(quantity)
		if err != nil {
			fmt.Printf("marshal: %v", err)
			os.Exit(1)
		}

		fmt.Println(string(bytes))
		return
	}

	if format == "csv" {
		return
	}
}
