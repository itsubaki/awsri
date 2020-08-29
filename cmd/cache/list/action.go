package list

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	dir := c.GlobalString("dir")

	files, err := walk(dir)
	if err != nil {
		fmt.Printf("dir walk: %v", err)
		os.Exit(1)
	}

	for _, f := range files {
		fmt.Println(f)
	}
}

func walk(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []string{}, err
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			files, err := walk(filepath.Join(dir, file.Name()))
			if err != nil {
				return []string{}, err
			}

			paths = append(paths, files...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths, nil
}
