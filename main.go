package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/cmd"
	"github.com/urfave/cli"
)

var date, hash, goversion string

func New(version string) *cli.App {
	app := cli.NewApp()

	app.Name = "hermes"
	app.Version = version
	app.Action = cmd.Action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Value: "stdout",
			Usage: "stdout",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "json",
			Usage: "json, csv",
		},
	}

	return app
}

func main() {
	version := fmt.Sprintf("%s %s %s", date, hash, goversion)
	hermes := New(version)
	if err := hermes.Run(os.Args); err != nil {
		panic(err)
	}
}
