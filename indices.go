package main

import (
	"context"
	"fmt"
	"io/ioutil"
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

	client, err := initClient(c.String("profile"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "initClient failure")
		return err
	}

	aliasService := elastic.NewAliasService(client)

	if c.Bool("rm") {
		_, err = aliasService.Remove(indexName, aliasName).Do(context.Background())
		if err != nil {
			return err
		}
	} else {
		_, err = aliasService.Add(indexName, aliasName).Do(context.Background())
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

	client, err := initClient(c.String("profile"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "initClient failure")
		return err
	}

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

	client, err := initClient(c.String("profile"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "initClient failure")
		return err
	}

	_, err = client.DeleteIndex(indexNames...).Do(context.Background())
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
