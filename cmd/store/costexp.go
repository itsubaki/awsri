package store

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/itsubaki/hermes/pkg/costexp"
	"github.com/urfave/cli"
)

func ActionStoreCostExp(c *cli.Context) {
	project := c.String("project")
	dir := c.GlobalString("dir")

	if err := StoreCostExp(project, dir); err != nil {
		fmt.Printf("store costexp: %v", err)
		os.Exit(1)
	}
}

func StoreCostExp(project, dir string) error {
	ctx := context.Background()
	ds, err := datastore.NewClient(ctx, project)
	if err != nil {
		return fmt.Errorf("new datastore client: %v", err)
	}

	kind := "costexp"
	date := costexp.GetCurrentDate()
	for i := range date {
		repo, err := costexp.Read(fmt.Sprintf("%s/%s/%s.out", dir, kind, date[i].YYYYMM()))
		if err != nil {
			return fmt.Errorf("read costexp: %v", err)
		}

		var key []*datastore.Key
		var src []interface{}

		rs := repo.SelectAll()
		for j := range rs {
			src = append(src, rs[j])
			key = append(key, &datastore.Key{
				Kind:   kind,
				Name:   rs[j].Hash(),
				Parent: nil,
			})
		}

		for j := range key {
			if _, err := ds.Put(ctx, key[j], src[j]); err != nil {
				return fmt.Errorf("put entity: %v", err)
			}
			fmt.Printf("put: %v\n", key[j])
		}
	}

	return nil
}
