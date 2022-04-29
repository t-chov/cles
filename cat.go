package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/olivere/elastic/v7"
	"github.com/urfave/cli/v2"
)

func prettyCatAliases(client *elastic.Client) (*string, error) {
	res, err := client.CatAliases().Human(true).Do(context.Background())
	if err != nil {
		return nil, err
	}

	var buf strings.Builder

	buf.WriteString(fmt.Sprintln("alias\tindex\trouting.index\trouting.search\tis_write_index\t"))
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

func cmdCatAliases(c *cli.Context) error {
	client, err := initClient(c.String("profile"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "initClient failure")
		return err
	}

	res, err := prettyCatAliases(client)
	if err != nil {
		return err
	}
	fmt.Print(*res)
	return nil
}
