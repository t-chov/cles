package main

import (
	"github.com/urfave/cli/v2"
)

func loadIndicesCommand() *cli.Command {
	return &cli.Command{
		Name:    "indices",
		Aliases: []string{"i", "index"},
		Usage:   "operate indices",
		Action:  cmdCatIndices,
		Subcommands: []*cli.Command{
			{
				Name:      "alias",
				Aliases:   []string{"a"},
				Usage:     "manage alias",
				Action:    cmdAliasIndex,
				ArgsUsage: "<INDEX_NAME> <ALIAS_NAME>",
				Flags: []cli.Flag{
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
					&cli.PathFlag{
						Name:        "body",
						Aliases:     []string{"b"},
						Usage:       "path to request body",
						DefaultText: "stdin",
					},
				},
			},
			{
				Name:      "delete",
				Aliases:   []string{"rm"},
				Usage:     "delete index",
				Action:    cmdDeleteIndex,
				ArgsUsage: "<INDEX_NAME>...",
			},
			{
				Name:      "mapping",
				Aliases:   []string{"m"},
				Usage:     "get mapping",
				Action:    cmdGetMapping,
				ArgsUsage: "<INDEX_NAME>...",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "all",
						Usage: "show all mappings",
					},
				},
			},
		},
	}
}

func loadCatCommand() *cli.Command {
	return &cli.Command{
		Name:    "cat",
		Aliases: []string{"c"},
		Usage:   "exec cat API",
		Subcommands: []*cli.Command{
			{
				Name:    "aliases",
				Aliases: []string{"a"},
				Usage:   "cat aliases",
				Action:  cmdCatAliases,
			},
			{
				Name:    "indices",
				Aliases: []string{"i"},
				Usage:   "cat indices",
				Action:  cmdCatIndices,
			},
		},
	}
}

func loadSearchTemplateCommand() *cli.Command {
	return &cli.Command{
		Name:    "search-template",
		Aliases: []string{"st"},
		Usage:   "operate search templates",
		Action:  cmdListSearchTemplates,
		Subcommands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "list search template",
				Action:  cmdListSearchTemplates,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
						Usage:   "show detail",
					},
				},
			},
			{
				Name:      "create",
				Aliases:   []string{"c", "new"},
				Usage:     "create search template",
				ArgsUsage: "<SEARCH_TEMPLATE_NAME>",
				Action:    cmdCreateSearchTemplate,
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:        "body",
						Aliases:     []string{"b"},
						Usage:       "path to request body",
						DefaultText: "stdin",
					},
				},
			},
			{
				Name:      "delete",
				Aliases:   []string{"rm"},
				Usage:     "delete search template",
				ArgsUsage: "<SEARCH_TEMPLATE_NAME>",
				Action:    cmdDeleteSearchTemplate,
			},
			{
				Name:      "render",
				Aliases:   []string{"r"},
				Usage:     "render search template",
				ArgsUsage: "<SEARCH_TEMPLATE_NAME>",
				Action:    cmdRenderSearchTemplate,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "params",
						Usage:       "set search params",
						DefaultText: "stdin",
					},
				},
			},
			{
				Name:      "search",
				Aliases:   []string{"s"},
				Usage:     "search with template",
				ArgsUsage: "<SEARCH_TEMPLATE_NAME>",
				Action:    cmdSearchWithTemplate,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "params",
						Usage:       "set search params",
						DefaultText: "stdin",
					},
					&cli.StringFlag{
						Name:    "index",
						Aliases: []string{"i"},
						Usage:   "index to search",
					},
				},
			},
		},
	}
}

func loadBulkCommand() *cli.Command {
	return &cli.Command{
		Name:    "bulk",
		Aliases: []string{"b"},
		Usage:   "operate bulk API",
		Subcommands: []*cli.Command{
			{
				Name:      "index",
				Aliases:   []string{"i"},
				Usage:     "exec bulk index from ndjson",
				ArgsUsage: "<INDEX_NAME>",
				Action:    cmdBulkIndex,
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:        "source",
						Aliases:     []string{"s", "src"},
						Usage:       "source file path(ndjson)",
						DefaultText: "stdin",
					},
					&cli.StringFlag{
						Name:    "id-column",
						Aliases: []string{"i"},
						Usage:   "column name for doc id",
					},
				},
			},
		},
	}
}

func loadReindexCommand() *cli.Command {
	return &cli.Command{
		Name:      "reindex",
		Aliases:   []string{"ri"},
		Usage:     "execute reindex API",
		Action:    cmdReindex,
		ArgsUsage: "<SOURCE> <DEST>",
	}
}

func loadCommands() []*cli.Command {
	commands := []*cli.Command{
		loadIndicesCommand(),
		loadCatCommand(),
		loadSearchTemplateCommand(),
		loadBulkCommand(),
		loadReindexCommand(),
	}
	return commands
}
