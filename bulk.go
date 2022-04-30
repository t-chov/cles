package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/olivere/elastic/v7"
	"github.com/urfave/cli/v2"
)

func id2str(rawId interface{}) (string, error) {
	switch val := rawId.(type) {
	case int:
		return strconv.Itoa(val), nil
	case string:
		return val, nil
	default:
		return "", fmt.Errorf("cannot convert id: %v", rawId)
	}
}

func cmdBulkIndex(c *cli.Context) error {
	client := c.Context.Value("client").(*elastic.Client)

	// set index name
	indexName := c.Args().First()
	if len(indexName) == 0 {
		return fmt.Errorf("need argument: name of index")
	}

	f, err := os.Open(c.Path("source"))
	if err != nil {
		return err
	}
	defer f.Close()

	idColumn := c.String("id-column")
	bulkRequest := client.Bulk().Index(indexName)
	// TODO make buffer size configurable. implement parse logic
	// default: 10MB
	bufferSize := 10 * 1024 * 1024

	var rowInterface interface{}
	currentBufferSize := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &rowInterface)
		if err != nil {
			return err
		}

		row := rowInterface.(map[string]interface{})
		if rawId, ok := row[idColumn]; ok {
			id, err := id2str(rawId)
			if err != nil {
				return err
			}
			request := elastic.NewBulkIndexRequest().Index(indexName).Id(id).Doc(row)
			estimatedBufferSize := currentBufferSize + len(request.String())
			if estimatedBufferSize > bufferSize {
				bulkRequest.Do(context.Background())
				os.Exit(0)
			} else {
				currentBufferSize = estimatedBufferSize
				bulkRequest.Add(request)
			}
		} else if len(idColumn) > 0 {
			return fmt.Errorf("no id column `%s` in %v", idColumn, row)
		}
	}
	bulkRequest.Do(context.Background())

	fmt.Printf("%v\n", client)

	return nil
}
