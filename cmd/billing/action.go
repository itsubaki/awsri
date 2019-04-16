package billing

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/pkg/billing"
	"github.com/itsubaki/hermes/pkg/costexp"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	dir := c.GlobalString("dir")
	date := costexp.GetCurrentDate()

	repo, err := NewRepository(dir, date)
	if err != nil {
		fmt.Println(fmt.Errorf("new billing repository: %v", err))
		os.Exit(1)
	}

	if c.String("format") == "csv" {
		// HEADER
		fmt.Print("AccountID,Description,")
		for i := len(date) - 1; i > -1; i-- {
			fmt.Printf("%s,", date[i].YYYYMM())
		}
		fmt.Println()

		for _, desc := range repo.Description() {
			list := repo.SelectAll().Description(desc)
			fmt.Printf("%s,%s,", list[0].AccountID, list[0].Description)

			for i := len(date) - 1; i > -1; i-- {
				if len(list.Date(date[i].Start)) == 0 {
					fmt.Printf("0.0,")
					continue
				}
				fmt.Printf("%s,", list.Date(date[i].Start)[0].UnblendedCost)
			}
			fmt.Println()
		}
		return
	}

	for _, r := range repo.SelectAll() {
		fmt.Printf("%#v\n", r)
	}
}

func NewRepository(dir string, date []*costexp.Date) (*billing.Repository, error) {
	out := &billing.Repository{}

	for _, d := range date {
		path := fmt.Sprintf("%s/billing/%s.out", dir, d.YYYYMM())
		repo, err := billing.Read(path)
		if err != nil {
			return nil, fmt.Errorf("read billing (path=%s): %v", path, err)
		}

		out.Internal = append(out.Internal, repo.Internal...)
	}

	return out, nil
}
