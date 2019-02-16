package store

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	p := c.String("project")

	ctx := context.Background()
	ds, err := datastore.NewClient(ctx, p)
	if err != nil {
		fmt.Println(fmt.Errorf("new datastore client: %v", err))
	}

	dir := c.GlobalString("dir")
	region := c.StringSlice("region")
	namespace := "pricing"

	for i := range region {
		repo, err := pricing.Read(fmt.Sprintf("%s/%s/%s.out", dir, namespace, region[i]))
		if err != nil {
			fmt.Println(fmt.Errorf("read pricing (region=%s): %v", region[i], err))
		}

		var key []*datastore.Key
		var src []interface{}

		rs := repo.SelectAll()
		for i := range rs {
			src = append(src, rs[i])
			key = append(key, &datastore.Key{
				Kind:      region[i],
				Name:      rs[i].ID(),
				Parent:    nil,
				Namespace: namespace,
			})
		}

		if _, err := ds.PutMulti(ctx, key, src); err != nil {
			fmt.Println(fmt.Errorf("put entity: %v", err))
		}
	}
}
