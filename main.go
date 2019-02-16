package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/cmd"
	"github.com/itsubaki/hermes/cmd/initialize"
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
		cli.StringFlag{
			Name:  "dir, d",
			Value: "/var/tmp/hermes",
		},
	}

	region := cli.StringSliceFlag{
		Name: "region, r",
		Value: &cli.StringSlice{
			"ap-northeast-1",
			"eu-central-1",
			"us-west-1",
			"us-west-2",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "init",
			Action: initialize.Action,
			Usage:  "download aws pricing, usage, reservation",
			Flags: []cli.Flag{
				region,
			},
		},
		{
			Name:   "store",
			Action: store.Action,
			Usage:  "import google datastore",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "project, p",
				},
				region,
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
