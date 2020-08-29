package rm

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/itsubaki/hermes/cmd/cache"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	dir := c.GlobalString("dir")
	yes := c.Bool("yes")

	files, err := cache.Dirwalk(dir)
	if err != nil {
		fmt.Printf("dir walk: %v", err)
		os.Exit(1)
	}

	if len(files) < 1 {
		fmt.Println("file not found")
		return
	}

	if yes {
		remove(files)
		return
	}

	for _, f := range files {
		fmt.Println(f)
	}

	fmt.Printf("confirm [y/N]: ")
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	yn := strings.ToLower(s.Text())

	if yn == "y" || yn == "yes" {
		remove(files)
		return
	}
}

func remove(files []string) {
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			fmt.Printf("remove file=%s: %v", f, err)
			os.Exit(1)
		}

		fmt.Printf("removed: %s\n", f)
	}
}
