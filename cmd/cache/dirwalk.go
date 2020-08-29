package cache

import (
	"io/ioutil"
	"path/filepath"
)

func Dirwalk(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []string{}, err
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			files, err := Dirwalk(filepath.Join(dir, file.Name()))
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
