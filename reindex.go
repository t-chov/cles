package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic/v7"
	"github.com/urfave/cli/v2"
)

func cmdReindex(c *cli.Context) error {
	if c.Args().Len() < 2 {
		return fmt.Errorf("1st arg is source index, 2nd arg is dest index. need two arguments")
	}
	src := c.Args().Get(0)
	dst := c.Args().Get(1)

	client := c.Context.Value("client").(*elastic.Client)
	svc := client.Reindex()
	res, err := svc.SourceIndex(src).DestinationIndex(dst).Do(context.TODO())
	if err != nil {
		return err
	}
	output, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))

	return nil
}
