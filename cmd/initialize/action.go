package initialize

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/pkg/costexp"
	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/reservation"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	region := c.StringSlice("region")
	dir := c.GlobalString("dir")

	Reservation(region, dir)
	Pricing(region, dir)
	CostExp(dir)
}

func Pricing(region []string, dir string) {
	path := fmt.Sprintf("%s/pricing", dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	for _, r := range region {
		cache := fmt.Sprintf("%s/%s.out", path, r)
		if _, err := os.Stat(cache); os.IsNotExist(err) {
			repo := pricing.NewRepository([]string{r})
			if err := repo.Fetch(); err != nil {
				fmt.Println(fmt.Errorf("fetch pricing (region=%s): %v", r, err))
				return
			}

			if err := repo.Write(cache); err != nil {
				fmt.Println(fmt.Errorf("write pricing (region=%s): %v", r, err))
				return
			}

			fmt.Printf("write: %v\n", cache)
		}
	}
}

func Reservation(region []string, dir string) {
	path := fmt.Sprintf("%s/reservation.out", dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		repo := reservation.NewRepository(region)
		if err := repo.Fetch(); err != nil {
			fmt.Println(fmt.Errorf("fetch reservation: %v", err))
			return
		}

		if err := repo.Write(path); err != nil {
			fmt.Println(fmt.Errorf("write reservation: %v", err))
			return
		}
	}

	fmt.Printf("write: %v\n", path)
}

func CostExp(dir string) {
	path := fmt.Sprintf("%s/costexp", dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	date := costexp.GetCurrentDate()
	for i := range date {
		cache := fmt.Sprintf("%s/%s.out", path, date[i].YYYYMM())
		if _, err := os.Stat(cache); os.IsNotExist(err) {
			repo := costexp.NewRepository([]*costexp.Date{date[i]})
			if err := repo.Fetch(); err != nil {
				fmt.Println(fmt.Errorf("fetch costexp (region=%s): %v", date[i], err))
				return
			}

			if err := repo.Write(cache); err != nil {
				fmt.Println(fmt.Errorf("write costexp (region=%s): %v", date[i], err))
				return
			}

			fmt.Printf("write: %v\n", cache)
		}
	}
}
