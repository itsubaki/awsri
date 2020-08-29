package list

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/cmd/cache"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	dir := c.GlobalString("dir")

	files, err := cache.Dirwalk(dir)
	if err != nil {
		fmt.Printf("dir walk: %v", err)
		os.Exit(1)
	}

	for _, f := range files {
		fmt.Println(f)
	}
}
