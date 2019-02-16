package initialize

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/pkg/costexp"
	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/reserved"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	region := c.StringSlice("region")
	dir := c.GlobalString("dir")

	if err := Pricing(region, dir); err != nil {
		fmt.Printf("write pricing: %v", err)
		os.Exit(1)
	}

	if err := CostExp(dir); err != nil {
		fmt.Printf("write costexp: %v", err)
		os.Exit(1)
	}

	if err := Reserved(region, dir); err != nil {
		fmt.Printf("write reservation: %v", err)
		os.Exit(1)
	}
}

func Pricing(region []string, dir string) error {
	path := fmt.Sprintf("%s/pricing", dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	for _, r := range region {
		cache := fmt.Sprintf("%s/%s.out", path, r)
		if _, err := os.Stat(cache); !os.IsNotExist(err) {
			continue
		}

		repo := pricing.NewRepository([]string{r})
		if err := repo.Fetch(); err != nil {
			return fmt.Errorf("fetch pricing (region=%s): %v", r, err)
		}

		if err := repo.Write(cache); err != nil {
			return fmt.Errorf("write pricing (region=%s): %v", r, err)
		}

		fmt.Printf("write: %v\n", cache)

	}

	return nil
}

func Reserved(region []string, dir string) error {
	path := fmt.Sprintf("%s/reserved.out", dir)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}

	repo := reserved.NewRepository(region)
	if err := repo.Fetch(); err != nil {
		return fmt.Errorf("fetch reservation: %v", err)
	}

	if err := repo.Write(path); err != nil {
		return fmt.Errorf("write reservation: %v", err)
	}

	fmt.Printf("write: %v\n", path)
	return nil
}

func CostExp(dir string) error {
	path := fmt.Sprintf("%s/costexp", dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	date := costexp.GetCurrentDate()
	for i := range date {
		cache := fmt.Sprintf("%s/%s.out", path, date[i].YYYYMM())
		if _, err := os.Stat(cache); !os.IsNotExist(err) {
			continue
		}

		repo := costexp.NewRepository([]*costexp.Date{date[i]})
		if err := repo.Fetch(); err != nil {
			return fmt.Errorf("fetch costexp (region=%s): %v", date[i], err)
		}

		if err := repo.Write(cache); err != nil {
			return fmt.Errorf("write costexp (region=%s): %v", date[i], err)
		}

		fmt.Printf("write: %v\n", cache)
	}

	return nil
}
