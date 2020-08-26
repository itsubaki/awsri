package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/cmd/cost"
	"github.com/itsubaki/hermes/cmd/fetch"
	"github.com/itsubaki/hermes/cmd/org"
	"github.com/itsubaki/hermes/cmd/pricing"
	"github.com/itsubaki/hermes/cmd/recommend"
	"github.com/itsubaki/hermes/cmd/reservation/reserved"
	"github.com/itsubaki/hermes/cmd/reservation/utilization"
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
		Name:  "region, r",
		Usage: "ap-east-1, ap-south-1, ap-northeast-1, ap-northeast-2, ap-northeast-3, ap-southeast-1, ap-southeast-2, eu-north-1, eu-west-1, eu-west-2, eu-west-3, eu-central-1, us-east-1, us-east-2, us-west-1, us-west-2, us-gov-east-1, us-gov-west-1, ca-central-1, sa-east-1, me-south-1",
		Value: &cli.StringSlice{
			"ap-northeast-1",
			"ap-southeast-1",
			"us-west-1",
			"us-west-2",
		},
	}

	format := cli.StringFlag{
		Name:  "format, f",
		Value: "json",
		Usage: "json, csv",
	}

	period := cli.StringFlag{
		Name:  "period, p",
		Value: "12m",
	}

	fetch := cli.Command{
		Name:    "fetch",
		Aliases: []string{"f"},
		Action:  fetch.Action,
		Usage:   "fetch aws pricing, usage, reservation, cost",
		Flags: []cli.Flag{
			region,
			period,
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

	cost := cli.Command{
		Name:    "cost",
		Aliases: []string{"c"},
		Action:  cost.Action,
		Usage:   "output cost group by linked account",
		Flags: []cli.Flag{
			format,
			period,
			cli.StringFlag{
				Name:  "attribute, a",
				Usage: "blended, unblended, net-unblended, amortized, net-amortized (format csv only)",
				Value: "unblended",
			},
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
			period,
			cli.BoolFlag{
				Name:  "normalize, n",
				Usage: "output normalized usage",
			},
			cli.BoolFlag{
				Name:  "merge, m",
				Usage: "output merged usage group by linked account",
			},
			cli.BoolFlag{
				Name:  "merge-overall, mm",
				Usage: "output merged usage",
			},
			cli.BoolFlag{
				Name:  "groupby, g",
				Usage: "output group by month/day usage (format json only)",
			},
			cli.StringFlag{
				Name:  "attribute, a",
				Usage: "hours, num (format csv only)",
				Value: "hours",
			},
		},
	}

	rsvutil := cli.Command{
		Name:    "utilization",
		Aliases: []string{"u"},
		Usage:   "output reservation utilization and coverage group by linked account",
		Action:  utilization.Action,
		Flags: []cli.Flag{
			region,
			format,
			period,
			cli.BoolFlag{
				Name:  "normalize, n",
				Usage: "output normalized usage",
			},
			cli.BoolFlag{
				Name:  "merge, m",
				Usage: "output merged usage group by linked account",
			},
			cli.BoolFlag{
				Name:  "groupby, g",
				Usage: "output group by month/day usage (format json only)",
			},
			cli.StringFlag{
				Name:  "attribute, a",
				Usage: "hours, num, percentage, ondemand-conversion-cost (format csv only)",
				Value: "hours",
			},
		},
	}

	rsved := cli.Command{
		Name:    "reserved",
		Aliases: []string{"r"},
		Usage:   "output history of reserved request",
		Action:  reserved.Action,
		Flags: []cli.Flag{
			region,
			format,
		},
	}

	reservation := cli.Command{
		Name:    "reservation",
		Aliases: []string{"r"},
		Usage:   "output reservation utilization, coverage, reserved",
		Subcommands: []cli.Command{
			rsvutil,
			rsved,
		},
	}

	recommend := cli.Command{
		Name:   "recommend",
		Action: recommend.Action,
		Usage:  "output recommended reserved instance num",
		Flags: []cli.Flag{
			format,
		},
	}

	org := cli.Command{
		Name:   "org",
		Action: org.Action,
		Usage:  "output list of accounts",
		Flags: []cli.Flag{
			format,
		},
	}

	app.Commands = []cli.Command{
		fetch,
		pricing,
		cost,
		usage,
		reservation,
		recommend,
		org,
	}

	return app
}

func main() {
	v := fmt.Sprintf("%s %s %s", date, hash, goversion)
	if err := New(v).Run(os.Args); err != nil {
		panic(err)
	}
}
