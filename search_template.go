package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/olivere/elastic/v7"
	"github.com/urfave/cli/v2"
)

func prettyListTemplate(client *elastic.Client, verbose bool) (*string, error) {
	state, err := client.ClusterState().Do(context.Background())
	if err != nil {
		return nil, err
	}

	templates := state.Metadata.StoredScripts
	var outputString string
	if verbose {
		output, err := json.MarshalIndent(templates, "", "  ")
		if err != nil {
			return nil, err
		}
		outputString = string(output)
	} else {
		var buf strings.Builder
		for key := range templates {
			buf.WriteString(fmt.Sprintln(key))
		}
		outputString = buf.String()
	}
	return &outputString, nil
}

func cmdListSearchTemplates(c *cli.Context) error {
	client, err := initClient(c.String("profile"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "initClient failure")
		return err
	}

	output, err := prettyListTemplate(client, c.Bool("verbose"))
	if err != nil {
		return err
	}

	fmt.Print(*output)
	return nil
}

func cmdCreateSearchTemplate(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return fmt.Errorf("must set search template name to create")
	}
	templateName := c.Args().Get(0)

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

	service := client.PutScript()
	_, err = service.Id(templateName).BodyString(string(bytes)).Do(context.Background())
	if err != nil {
		return err
	}

	output, err := prettyListTemplate(client, c.Bool("verbose"))
	if err != nil {
		return err
	}

	fmt.Print(*output)
	return nil
}

func cmdDeleteSearchTemplate(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return fmt.Errorf("must set search template name to create")
	}
	templateName := c.Args().Get(0)

	client, err := initClient(c.String("profile"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "initClient failure")
		return err
	}

	service := client.DeleteScript()
	_, err = service.Id(templateName).Do(context.Background())
	if err != nil {
		return err
	}

	output, err := prettyListTemplate(client, c.Bool("verbose"))
	if err != nil {
		return err
	}

	fmt.Print(*output)
	return nil
}
