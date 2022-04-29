package main

import (
	"context"
	"fmt"
	"os"

	"github.com/olivere/elastic/v7"
	"github.com/urfave/cli/v2"
)

func cmdAliasIndex(c *cli.Context) error {
	if c.Args().Len() < 2 {
		return fmt.Errorf("1st arg is index, 2nd arg ais aliases. need two arguments")
	}
	indexName := c.Args().Get(0)
	aliasName := c.Args().Get(1)

	fmt.Printf("index: %s, alias: %s\n", indexName, aliasName)

	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "client failure")
		return err
	}

	aliasService := elastic.NewAliasService(client)

	res, err := aliasService.Add(indexName, aliasName).Do(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", res)
	return nil
}
