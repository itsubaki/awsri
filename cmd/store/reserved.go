package store

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/itsubaki/hermes/pkg/reserved"
	"github.com/urfave/cli"
)

func ActionStoreReserved(c *cli.Context) {
	project := c.String("project")
	dir := c.GlobalString("dir")

	if err := StoreReserved(project, dir); err != nil {
		fmt.Printf("store reserved: %v", err)
		os.Exit(1)
	}
}

func StoreReserved(project, dir string) error {
	ctx := context.Background()
	ds, err := datastore.NewClient(ctx, project)
	if err != nil {
		return fmt.Errorf("new datastore client: %v", err)
	}

	kind := "reserved"
	repo, err := reserved.Read(fmt.Sprintf("%s/%s.out", dir, kind))
	if err != nil {
		return fmt.Errorf("read reservation: %v", err)
	}

	var key []*datastore.Key
	var src []interface{}

	rs := repo.SelectAll()
	for j := range rs {
		src = append(src, rs[j])
		key = append(key, &datastore.Key{
			Kind:   kind,
			Name:   rs[j].ReservedID,
			Parent: nil,
		})
	}

	for j := range key {
		if _, err := ds.Put(ctx, key[j], src[j]); err != nil {
			return fmt.Errorf("put entity: %v", err)
		}
		fmt.Printf("put: %v\n", key[j])
	}

	return nil
}
