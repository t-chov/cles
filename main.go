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
						Aliases:     []string{"p"},
						Usage:       "set profile name",
						DefaultText: "default",
					},
					&cli.BoolFlag{
						Name:    "delete",
						Aliases: []string{"rm", "d"},
						Usage:   "delete alias",
					},
				},
			},
			{
				Name:      "create",
				Aliases:   []string{"c", "new"},
				Usage:     "create index",
				Action:    cmdCreateIndex,
				ArgsUsage: "<INDEX_NAME>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "profile",
						Aliases:     []string{"p"},
						Usage:       "set profile name",
						DefaultText: "default",
					},
					&cli.PathFlag{
						Name:     "body",
						Aliases:  []string{"b"},
						Usage:    "path to request body",
						Required: true,
					},
				},
			},
			{
				Name:      "delete",
				Aliases:   []string{"d", "rm"},
				Usage:     "delete index",
				Action:    cmdDeleteIndex,
				ArgsUsage: "<INDEX_NAME>...",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "profile",
						Aliases:     []string{"p"},
						Usage:       "set profile name",
						DefaultText: "default",
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
						Aliases:     []string{"p"},
						Usage:       "set profile name",
						DefaultText: "default",
					},
				},
			},
			{
				Name:    "indices",
				Aliases: []string{"i"},
				Usage:   "cat indices",
				Action:  cmdCatIndices,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "profile",
						Aliases:     []string{"p"},
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
