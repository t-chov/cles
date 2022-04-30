package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/olivere/elastic/v7"
	"github.com/urfave/cli/v2"
)

func cmdAliasIndex(c *cli.Context) error {
	if c.Args().Len() < 2 {
		return fmt.Errorf("1st arg is index, 2nd arg ais aliases. need two arguments")
	}
	indexName := c.Args().Get(0)
	aliasName := c.Args().Get(1)

	client := c.Context.Value("client").(*elastic.Client)

	aliasService := elastic.NewAliasService(client)

	if c.Bool("rm") {
		_, err := aliasService.Remove(indexName, aliasName).Do(context.Background())
		if err != nil {
			return err
		}
	} else {
		_, err := aliasService.Add(indexName, aliasName).Do(context.Background())
		if err != nil {
			return err
		}
	}

	res, err := prettyCatAliases(client)
	if err != nil {
		return err
	}

	fmt.Print(*res)
	return nil
}

func cmdCreateIndex(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return fmt.Errorf("must set index name to create")
	}
	indexName := c.Args().Get(0)

	client := c.Context.Value("client").(*elastic.Client)

	body := c.Path("body")
	bytes, err := ioutil.ReadFile(body)
	if err != nil {
		return err
	}

	service := client.CreateIndex(indexName)
	_, err = service.Body(string(bytes)).Do(context.Background())
	if err != nil {
		return err
	}

	res, err := prettyCatIndices(client)
	if err != nil {
		return err
	}
	fmt.Print(*res)
	return nil
}

func cmdDeleteIndex(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return fmt.Errorf("must set index names to delete")
	}
	indexNames := c.Args().Slice()

	client := c.Context.Value("client").(*elastic.Client)

	_, err := client.DeleteIndex(indexNames...).Do(context.Background())
	if err != nil {
		return err
	}

	res, err := prettyCatIndices(client)
	if err != nil {
		return err
	}

	fmt.Print(*res)
	return nil
}

func cmdGetMapping(c *cli.Context) error {
	index := c.Args().Slice()
	if c.Bool("all") {
		index = []string{"_all"}
	}

	client := c.Context.Value("client").(*elastic.Client)

	service := client.GetMapping()
	res, err := service.Index(index...).Do(context.Background())
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))

	return nil
}
