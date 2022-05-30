package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/olivere/elastic/v7"
	"github.com/urfave/cli/v2"
)

type catAliasesSvc interface {
	Do(ctx context.Context) (elastic.CatAliasesResponse, error)
}

type catIndicesSvc interface {
	Do(ctx context.Context) (elastic.CatIndicesResponse, error)
}

func prettyCatAliases(service catAliasesSvc) (*string, error) {
	res, err := service.Do(context.Background())
	if err != nil {
		return nil, err
	}

	var buf strings.Builder

	buf.WriteString(fmt.Sprintln("alias\tindex\trouting.index\trouting.search\tis_write_index"))
	for _, v := range res {
		row := fmt.Sprintf(
			"%s\t%s\t%s\t%s\t%s\n",
			v.Alias,
			v.Index,
			v.RoutingIndex,
			v.RoutingSearch,
			v.IsWriteIndex,
		)
		buf.WriteString(row)
	}
	output := buf.String()
	return &output, nil
}

func prettyCatIndices(service catIndicesSvc) (*string, error) {
	res, err := service.Do(context.Background())
	if err != nil {
		return nil, err
	}

	var buf strings.Builder
	headers := []string{
		"health", "status", "index", "uuid", "pri", "rep", "docs.count",
		"docs.deleted", "store.size", "pri.store.size",
	}
	buf.WriteString(strings.Join(headers, "\t") + "\n")
	for _, v := range res {
		row := fmt.Sprintf(
			"%s\t%s\t%s\t%s\t%d\t%d\t%d\t%d\t%s\t%s\n",
			v.Health, v.Status, v.Index, v.UUID, v.Pri, v.Rep, v.DocsCount,
			v.DocsDeleted, v.StoreSize, v.PriStoreSize,
		)
		buf.WriteString(row)
	}
	output := buf.String()
	return &output, nil
}

func cmdCatAliases(c *cli.Context) error {
	client := c.Context.Value("client").(*elastic.Client)

	res, err := prettyCatAliases(client.CatAliases().Human(true))
	if err != nil {
		return err
	}
	fmt.Print(*res)
	return nil
}

func cmdCatIndices(c *cli.Context) error {
	client := c.Context.Value("client").(*elastic.Client)

	res, err := prettyCatIndices(client.CatIndices().Human(true))
	if err != nil {
		return err
	}
	fmt.Print(*res)
	return nil
}
