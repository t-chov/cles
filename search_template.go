package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/olivere/elastic/v7"
	"github.com/urfave/cli/v2"
)

func buildTemplateParameters(templateId string, rawParams string, stream *os.File) (map[string]interface{}, error) {
	var params interface{}
	if len(rawParams) == 0 {
		bytes, err := ioutil.ReadAll(stream)
		if err != nil {
			return nil, err
		}
		rawParams = string(bytes)
	}
	err := json.Unmarshal([]byte(rawParams), &params)
	if err != nil {
		return nil, err
	}
	searchBody := make(map[string]interface{})
	searchBody["id"] = templateId
	searchBody["params"] = params
	return searchBody, nil
}

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
	var bytes []byte
	var err error
	if len(body) > 0 {
		bytes, err = ioutil.ReadFile(body)
	} else {
		bytes, err = ioutil.ReadAll(os.Stdin)
	}
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

	renderBody, err := buildTemplateParameters(templateName, c.String("params"), os.Stdin)
	if err != nil {
		return err
	}

	var urlValues url.Values
	res, err := client.PerformRequest(context.Background(), elastic.PerformRequestOptions{
		Method:      "POST",
		Path:        "/_render/template",
		Params:      urlValues,
		Body:        renderBody,
		ContentType: "application/json",
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", res.Body)
	return nil
}

func cmdSearchWithTemplate(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return fmt.Errorf("must set search template name to create")
	}
	templateName := c.Args().First()

	client := c.Context.Value("client").(*elastic.Client)

	searchBody, err := buildTemplateParameters(templateName, c.String("params"), os.Stdin)
	if err != nil {
		return err
	}

	var path string
	index := c.String("index")
	if len(index) > 0 {
		path = fmt.Sprintf("/%s/_search/template", index)
	} else {
		path = "/_search/template"
	}

	var urlValues url.Values
	res, err := client.PerformRequest(context.Background(), elastic.PerformRequestOptions{
		Method:      "POST",
		Path:        path,
		Params:      urlValues,
		Body:        searchBody,
		ContentType: "application/json",
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", res.Body)
	return nil
}
