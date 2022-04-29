package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli/v2"
)

func cmdListSearchTemplates(c *cli.Context) error {
	client, err := initClient(c.String("profile"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "initClient failure")
		return err
	}

	state, err := client.ClusterState().Do(context.Background())
	if err != nil {
		return err
	}

	templates := state.Metadata.StoredScripts
	if c.Bool("verbose") {
		output, err := json.MarshalIndent(templates, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
	} else {
		for key := range templates {
			fmt.Println(key)
		}
	}

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

	return nil
}
