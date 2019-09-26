package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/cmd/fetch"
	"github.com/itsubaki/hermes/cmd/pricing"
	"github.com/urfave/cli"
)

var date, hash, goversion string

func New(version string) *cli.App {
	app := cli.NewApp()

	app.Name = "hermes"
	app.Usage = "aws cost optimization"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "dir, d",
			Value: "/var/tmp/hermes",
		},
	}

	region := cli.StringSliceFlag{
		Name: "region, r",
		Value: &cli.StringSlice{
			"ap-northeast-1",
			"us-west-2",
		},
	}

	format := cli.StringFlag{
		Name:  "format, f",
		Value: "json",
		Usage: "json, csv",
	}

	fetch := cli.Command{
		Name:    "fetch",
		Aliases: []string{"f"},
		Action:  fetch.Action,
		Usage:   "fetch aws pricing, usage",
		Flags: []cli.Flag{
			region,
		},
	}

	pricing := cli.Command{
		Name:    "pricing",
		Aliases: []string{"p"},
		Usage:   "output aws pricing",
		Action:  pricing.Action,
		Flags: []cli.Flag{
			region,
			format,
		},
	}

	usage := cli.Command{
		Name:    "usage",
		Aliases: []string{"u"},
		Usage:   "output aws instance hour usage",
		Flags: []cli.Flag{
			region,
			format,
		},
	}

	app.Commands = []cli.Command{
		fetch,
		pricing,
		usage,
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
