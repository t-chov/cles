package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	appName  = "cles"
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
				ArgsUsage: "<INDEX_NAME> <ALIAS_NAME>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "profile",
						Usage:       "set profile name",
						DefaultText: "default",
					},
					&cli.BoolFlag{
						Name:  "rm",
						Usage: "remove alias",
					},
				},
			},
		},
	},
	{
		Name:    "cat",
		Aliases: []string{"c"},
		Usage:   "exec cat API",
		Subcommands: []*cli.Command{
			{
				Name:    "aliases",
				Aliases: []string{"a"},
				Usage:   "cat aliases",
				Action:  cmdCatAliases,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "profile",
						Usage:       "set profile name",
						DefaultText: "default",
					},
				},
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
	app.Name = appName
	app.Usage = "Command line client for Elasticsearch"
	app.Version = version
	app.Commands = commands
	app.Action = appRun

	return msg(app.Run(os.Args))
}

func main() {
	os.Exit(run())
}
