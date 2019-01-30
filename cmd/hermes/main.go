package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var date, hash, goversion string

func main() {
	version := fmt.Sprintf("%s %s %s", date, hash, goversion)
	hermes := New(version)
	if err := hermes.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func New(version string) *cli.App {
	app := cli.NewApp()

	app.Name = "hermes"
	app.Version = version
	app.Commands = []cli.Command{}

	return app
}
