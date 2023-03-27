package main

import (
	"context"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

const (
	debugColor = color.FgCyan
	errorColor = color.FgRed
)

var commands = loadCommands()

func setClient(c *cli.Context) error {
	debugFn := c.Context.Value("debugFunc").(func(message string))
	client, err := initClient(c.String("profile"), debugFn)
	if err != nil {
		return err
	}
	//lint:ignore SA1029 initClient before subcommand
	c.Context = context.WithValue(c.Context, "client", client)
	return nil
}

func msg(err error) int {
	if err != nil {
		red := color.New(errorColor).FprintfFunc()
		red(os.Stderr, "%s: %v\n", os.Args[0], err)
		return 1
	}
	return 0
}

func run(app *cli.App) int {
	return msg(app.Run(os.Args))
}

func main() {
	os.Exit(run(initApp()))
}
