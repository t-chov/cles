package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/cheggaaa/pb/v3"
	"github.com/olivere/elastic/v7"
	"github.com/urfave/cli/v2"
)

func idToStr(rawId interface{}) (string, error) {
	switch val := rawId.(type) {
	case int:
		return strconv.Itoa(val), nil
	case string:
		return val, nil
	default:
		return "", fmt.Errorf("cannot convert id: %v", rawId)
	}
}

func calcNumOfLines(f *os.File) int {
	sc := bufio.NewScanner(f)
	num := 0
	for sc.Scan() {
		num++
	}
	f.Seek(0, io.SeekStart)
	return num
}

func rowToDoc(rawRow []byte) (map[string]interface{}, error) {
	var rowInterface interface{}
	err := json.Unmarshal(rawRow, &rowInterface)
	if err != nil {
		return nil, err
	}
	return rowInterface.(map[string]interface{}), nil
}

func getId(doc map[string]interface{}, idColumn string) (*string, error) {
	if len(idColumn) == 0 {
		return nil, nil
	}
	if rawId, ok := doc[idColumn]; ok {
		id, err := idToStr(rawId)
		if err != nil {
			return nil, err
		}
		return &id, nil
	} else {
		return nil, fmt.Errorf("no id column `%s` in %v", idColumn, doc)
	}
}

func cmdBulkIndex(c *cli.Context) error {
	client := c.Context.Value("client").(*elastic.Client)

	// set index name
	indexName := c.Args().First()
	if len(indexName) == 0 {
		return fmt.Errorf("need argument: name of index")
	}

	// load source file
	var f *os.File
	var bar *pb.ProgressBar
	if len(c.Path("source")) > 0 {
		var err error
		f, err = os.Open(c.Path("source"))
		if err != nil {
			return err
		}
		bar = pb.StartNew(calcNumOfLines(f))
	} else {
		f = os.Stdin
		bar = pb.StartNew(-1)
	}
	defer f.Close()

	bulkRequest := client.Bulk().Index(indexName)
	// TODO make buffer size configurable. implement parse logic
	// default: 10MB
	bufferSize := 10 * 1024 * 1024

	currentBufferSize := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		bar.Increment()
		doc, err := rowToDoc(scanner.Bytes())
		if err != nil {
			return err
		}

		id, err := getId(doc, c.String("id-column"))
		if err != nil {
			return err
		}

		eachRequet := elastic.NewBulkIndexRequest().Index(indexName).Doc(doc)
		if id != nil {
			eachRequet = eachRequet.Id(*id)
		}

		eachReqSize := len(eachRequet.String())
		estimatedBufferSize := currentBufferSize + eachReqSize
		if estimatedBufferSize > bufferSize {
			_, err := bulkRequest.Do(context.Background())
			if err != nil {
				return err
			}
			currentBufferSize = 0
		}
		currentBufferSize += eachReqSize
		bulkRequest.Add(eachRequet)
	}
	_, err := bulkRequest.Do(context.Background())
	if err != nil {
		return err
	}
	bar.Finish()

	return nil
}
