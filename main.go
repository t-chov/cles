package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	name     = "cles"
	version  = "0.0.1"
	revision = "HEAD"
)

var commands = []*cli.Command{
	{
		Name:    "indices",
		Aliases: []string{"i", "index"},
		Usage:   "operate indices",
		Subcommands: []*cli.Command{
			{
				Name:      "alias",
				Aliases:   []string{"a"},
				Usage:     "manage alias",
				Action:    cmdAliasIndex,
				ArgsUsage: "<INDEC_NAME> <ALIAS_NAMES>...",
			},
		},
	},
}

func msg(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return 1
	}
	return 0
}

func appRun(c *cli.Context) error {
	args := c.Args()
	if !args.Present() {
		cli.ShowAppHelp(c)
	}
	return nil
}

func run() int {
	app := cli.NewApp()
	app.Name = name
	app.Usage = "Command line client for Elasticsearch"
	app.Version = version
	app.Commands = commands
	app.Action = appRun

	return msg(app.Run(os.Args))
}

func main() {
	os.Exit(run())
}
