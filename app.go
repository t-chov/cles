package main

import (
	"context"

	"github.com/urfave/cli/v2"
)

const (
	APP_NAME = "cles"
	VERSION  = "0.0.8"
	REVISION = "HEAD"
)

func initializeContext(c *cli.Context) error {
	fp, err := initDebugFunc(c.Bool("debug"))
	if err != nil {
		return err
	}
	//lint:ignore SA1029 initDebugFunc before subcommand
	c.Context = context.WithValue(c.Context, "debugFunc", *fp)
	err = setClient(c)
	return err
}

func appRun(c *cli.Context) error {
	args := c.Args()
	if !args.Present() {
		cli.ShowAppHelp(c)
	}
	return nil
}

func initApp() *cli.App {
	app := cli.NewApp()
	app.Name = APP_NAME
	app.Usage = "Command line client for Elasticsearch"
	app.Version = VERSION
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
	return app
}
