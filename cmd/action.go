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
	price, err := NewPricingRepository(merged.Region(), dir)
	if err != nil {
		fmt.Println(fmt.Errorf("new pricing repository: %v", err))
		os.Exit(1)
	}

	rec, err := merged.Recommend(price)
	if err != nil {
		fmt.Println(fmt.Errorf("recommend: %v", err))
		os.Exit(1)
	}

	rsv, err := NewReservedRepository(dir)
	if err != nil {
		fmt.Println(fmt.Errorf("new reserved repository: %v", err))
		os.Exit(1)
	}

	cov, err := api.GetCoverage(rec.NormalizedList(), rsv)
	if err != nil {
		fmt.Println(fmt.Errorf("new coverage list: %v", err))
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

/*
NewPricingRepository returns list of pricing repository.
*/
func NewPricingRepository(region []string, dir string) ([]*pricing.Repository, error) {
	out := []*pricing.Repository{}
	for _, r := range region {
		path := fmt.Sprintf("%s/pricing/%s.out", dir, r)
		repo, err := pricing.Read(path)
		if err != nil {
			return nil, fmt.Errorf("read pricing (path=%s): %v", path, err)
		}
		out = append(out, repo)
	}

	return out, nil
}

/*
NewReservedRepository returns reserved repository.
*/
func NewReservedRepository(dir string) (*reserved.Repository, error) {
	path := fmt.Sprintf("%s/reserved.out", dir)
	repo, err := reserved.Read(path)
	if err != nil {
		return nil, fmt.Errorf("read reservation (path=%s): %v", path, err)
	}

	return repo, nil
}
