package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
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
	client := c.Context.Value("client").(*elastic.Client)

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
	templateName := c.Args().First()

	client := c.Context.Value("client").(*elastic.Client)

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
	templateName := c.Args().First()

	client := c.Context.Value("client").(*elastic.Client)

	service := client.DeleteScript()
	_, err := service.Id(templateName).Do(context.Background())
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

func cmdRenderSearchTemplate(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return fmt.Errorf("must set search template name to create")
	}
	templateName := c.Args().First()

	client := c.Context.Value("client").(*elastic.Client)

	var params interface{}
	paramsRaw := c.String("params")
	err := json.Unmarshal([]byte(paramsRaw), &params)
	if err != nil {
		return err
	}
	renderbody := make(map[string]interface{})
	renderbody["id"] = templateName
	renderbody["params"] = params

	var urlValues url.Values
	res, err := client.PerformRequest(context.Background(), elastic.PerformRequestOptions{
		Method:      "POST",
		Path:        "/_render/template",
		Params:      urlValues,
		Body:        renderbody,
		ContentType: "application/json",
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", res.Body)
	return nil
}
