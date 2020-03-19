package main

import (
	"fmt"
	"os"

	hermes "github.com/itsubaki/hermes/pkg"
)

var date, hash, goversion string

func main() {
	v := fmt.Sprintf("%s %s %s", date, hash, goversion)
	if err := hermes.New(v).Run(os.Args); err != nil {
		panic(err)
	}
}
