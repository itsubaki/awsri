package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/cmd"
	"github.com/itsubaki/hermes/cmd/store"
	"github.com/urfave/cli"
)

var date, hash, goversion string

func New(version string) *cli.App {
	app := cli.NewApp()

	app.Name = "hermes"
	app.Usage = "recommend aws reserved instances"
	app.Version = version
	app.Action = cmd.Action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "format, f",
			Value: "json",
			Usage: "json, csv, tsv",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "store",
			Action: store.Action,
			Usage:  "import google datastore",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "project, p",
				},
			},
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
