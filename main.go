package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/cmd"
)

var date, hash, goversion string

func main() {
	v := fmt.Sprintf("%s %s %s", date, hash, goversion)
	if err := cmd.New(v).Run(os.Args); err != nil {
		panic(err)
	}
}
