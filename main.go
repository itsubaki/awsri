package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/cmd"
	"github.com/itsubaki/hermes/cmd/initialize"
	"github.com/itsubaki/hermes/cmd/store/costexp"
	"github.com/itsubaki/hermes/cmd/store/pricing"
	"github.com/itsubaki/hermes/cmd/store/reserved"
	"github.com/urfave/cli"
)

var date, hash, goversion string

func New(version string) *cli.App {
	app := cli.NewApp()

	app.Name = "hermes"
	app.Usage = "aws cost optimization"
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
			"ap-southeast-1",
			"eu-central-1",
			"us-west-1",
			"us-west-2",
		},
	}

	init := cli.Command{
		Name:   "init",
		Action: initialize.Action,
		Usage:  "download aws pricing, usage, reservation",
		Flags: []cli.Flag{
			region,
		},
	}

	flags := []cli.Flag{
		cli.StringFlag{
			Name: "project, p",
		},
	}

	store := cli.Command{
		Name: "store",
		Subcommands: []cli.Command{
			{
				Name:    "pricing",
				Aliases: []string{"p"},
				Action:  pricing.Action,
				Flags:   append(flags, region),
			},
			{
				Name:    "costexp",
				Aliases: []string{"c"},
				Action:  costexp.Action,
				Flags:   flags,
			},
			{
				Name:    "reserved",
				Aliases: []string{"r"},
				Action:  reserved.Action,
				Flags:   flags,
			},
		},
	}

	app.Commands = []cli.Command{
		init,
		store,
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
