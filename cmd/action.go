package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/itsubaki/hermes/pkg/api"
	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/reserved"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(fmt.Errorf("stdin: %v", err))
		return
	}

	var input api.Input
	if uerr := json.Unmarshal(stdin, &input); uerr != nil {
		fmt.Println(fmt.Errorf("unmarshal: %v", uerr))
		os.Exit(1)
	}

	dir := c.GlobalString("dir")
	merged := input.Forecast.Merge()

	price := []*pricing.Repository{}
	for _, in := range merged {
		path := fmt.Sprintf("%s/pricing/%s.out", dir, in.Region)
		repo, rerr := pricing.Read(path)
		if rerr != nil {
			fmt.Println(fmt.Errorf("read pricing (region=%s): %v", in.Region, rerr))
			os.Exit(1)
		}
		price = append(price, repo)
	}

	rec, err := merged.Recommend(price)
	if err != nil {
		fmt.Println(fmt.Errorf("recommend: %v", err))
		os.Exit(1)
	}

	path := fmt.Sprintf("%s/reserved.out", dir)
	rsv, err := reserved.Read(path)
	if err != nil {
		fmt.Println(fmt.Errorf("read reservation: %v", err))
		os.Exit(1)
	}

	cov, err := api.GetCoverage(rec.Merge(), rsv)
	if err != nil {
		fmt.Println(fmt.Errorf("new result list: %v", err))
		os.Exit(1)
	}

	output := &api.Output{
		Forecast:    input.Forecast,
		Merged:      merged,
		Recommended: rec,
		Coverage:    cov,
	}

	if c.String("format") == "csv" {
		fmt.Println(output.CSV())
		return
	}

	if c.String("format") == "tsv" {
		fmt.Println(output.TSV())
		return
	}

	//  c.String("format") == "json"
	fmt.Println(output.JSON())
}
