package list

import (
	"fmt"

	"github.com/itsubaki/hermes/cmd/cache"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	dir := c.String("dir")

	files, err := cache.Dirwalk(dir)
	if err != nil {
		return fmt.Errorf("dir walk: %v", err)
	}

	for _, f := range files {
		fmt.Println(f)
	}

	return nil
}
