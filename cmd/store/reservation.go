package store

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/itsubaki/hermes/pkg/reservation"
	"github.com/urfave/cli"
)

func ActionStoreReservation(c *cli.Context) {
	project := c.String("project")
	dir := c.GlobalString("dir")

	if err := StoreReservation(project, dir); err != nil {
		fmt.Printf("store reservation: %v", err)
		os.Exit(1)
	}
}

func StoreReservation(project, dir string) error {
	ctx := context.Background()
	ds, err := datastore.NewClient(ctx, project)
	if err != nil {
		return fmt.Errorf("new datastore client: %v", err)
	}

	kind := "reservation"
	repo, err := reservation.Read(fmt.Sprintf("%s/%s.out", dir, kind))
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
