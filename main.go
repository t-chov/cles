package main

import (
	"context"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

const (
	appName    = "cles"
	version    = "0.0.7"
	revision   = "HEAD"
	debugColor = color.FgCyan
	errorColor = color.FgRed
)

var commands = loadCommands()

func debug(stream *os.File, message string) {
	debugFunc := color.New(debugColor).FprintfFunc()
	debugFunc(stream, message)
}

func setColoredWriter(c *cli.Context) error {
	if c.Bool("debug") {
		//lint:ignore SA1029 set debug stream
		c.Context = context.WithValue(c.Context, "debugStream", os.Stdout)
	} else {
		devNull, err := os.OpenFile(os.DevNull, os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		//lint:ignore SA1029 set debug stream
		c.Context = context.WithValue(c.Context, "debugStream", devNull)
	}
	return nil
}

func setClient(c *cli.Context) error {
	debugStream := c.Context.Value("debugStream").(*os.File)
	client, err := initClient(c.String("profile"), debugStream)
	if err != nil {
		return err
	}
	//lint:ignore SA1029 initClient before subcommand
	c.Context = context.WithValue(c.Context, "client", client)
	return nil
}

func initializeContext(c *cli.Context) error {
	setColoredWriter(c)
	err := setClient(c)
	return err
}

func msg(err error) int {
	if err != nil {
		red := color.New(errorColor).FprintfFunc()
		red(os.Stderr, "%s: %v\n", os.Args[0], err)
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
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "profile",
			Value:   "default",
			Aliases: []string{"p"},
			Usage:   "set profile name",
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "show detail log",
		},
	}
	app.Before = initializeContext

	return msg(app.Run(os.Args))
}

func main() {
	os.Exit(run())
}
