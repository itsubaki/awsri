package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/hermes/cmd/hermes/predict"
	"github.com/itsubaki/hermes/cmd/hermes/recommend"
	"github.com/urfave/cli"
)

var date, hash, goversion string

func New(version string) *cli.App {
	app := cli.NewApp()

	app.Name = "hermes"
	app.Version = version

	recommend := cli.Command{
		Name:    "recommend",
		Aliases: []string{"r"},
		Action:  recommend.Action,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "output, o",
				Value: "stdout",
				Usage: "stdout, googless",
			},
		},
	}

	predict := cli.Command{
		Name:    "predict",
		Aliases: []string{"p"},
		Action:  predict.Action,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "method, m",
				Value: "linear",
			},
		},
	}

	app.Commands = []cli.Command{
		recommend,
		predict,
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
