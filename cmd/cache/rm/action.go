package rm

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/itsubaki/hermes/cmd/cache"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	dir := c.String("dir")
	yes := c.Bool("yes")

	files, err := cache.Dirwalk(dir)
	if err != nil {
		return fmt.Errorf("dir walk: %v", err)
	}

	if len(files) < 1 {
		fmt.Println("file not found")
		return nil
	}

	if yes {
		remove(files)
		return nil
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
		return nil
	}

	return nil
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
