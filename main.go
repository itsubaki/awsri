package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/cmd/fetch"
	"github.com/itsubaki/hermes/cmd/pricing"
	"github.com/itsubaki/hermes/cmd/usage"
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
			"ap-northeast-2",
			"ap-northeast-3",
			"ap-southeast-1",
			"ap-southeast-2",
			"eu-west-1",
			"us-east-1",
			"us-east-2",
			"us-west-1",
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
		Action:  pricing.Action,
		Usage:   "output aws pricing",
		Flags: []cli.Flag{
			region,
			format,
		},
	}

	usage := cli.Command{
		Name:    "usage",
		Aliases: []string{"u"},
		Action:  usage.Action,
		Usage:   "output aws instance hour usage",
		Flags: []cli.Flag{
			region,
			format,
			cli.BoolFlag{
				Name:  "normalize, n",
				Usage: "output normalized aws instance hour usage",
			},
			cli.BoolFlag{
				Name:  "merge, m",
				Usage: "output merged aws instance hour usage",
			},
			cli.BoolFlag{
				Name:  "monthly, mm",
				Usage: "output monthly aws instance hour usage",
			},
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
