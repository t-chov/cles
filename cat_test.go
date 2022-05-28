package main

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/olivere/elastic/v7"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type catAliasesImpl struct {
	ExpectError bool
}

func (c catAliasesImpl) Do(ctx context.Context) (elastic.CatAliasesResponse, error) {
	if c.ExpectError {
		return []elastic.CatAliasesResponseRow{}, fmt.Errorf("cat aliases error")
	}
	resp := []elastic.CatAliasesResponseRow{
		{
			Alias:         "foo",
			Index:         "foo01",
			Filter:        "*",
			RoutingIndex:  "FooRoutingIndex",
			RoutingSearch: "FooRoutingSearch",
			IsWriteIndex:  "FooIsWriteIndex",
		},
		{
			Alias:         "bar",
			Index:         "bar02",
			Filter:        "*",
			RoutingIndex:  "BarRoutingIndex",
			RoutingSearch: "BarRoutingSearch",
			IsWriteIndex:  "BarIsWriteIndex",
		},
	}
	return resp, nil
}

func TestPrettyCatAliases(t *testing.T) {
	t.Run(
		"fail catAliases API",
		func(t *testing.T) {
			svc := catAliasesImpl{
				ExpectError: true,
			}
			res, err := prettyCatAliases(svc)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
			if !strings.Contains(fmt.Sprintf("%s", err), "cat aliases error") {
				t.Errorf("unexpected error: %v", err)
			}
			if res != nil {
				t.Errorf("expected nil, got %v", res)
			}
		},
	)
	t.Run(
		"success catAliases API",
		func(t *testing.T) {
			svc := catAliasesImpl{
				ExpectError: false,
			}
			res, err := prettyCatAliases(svc)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if res == nil {
				t.Errorf("expected not nil, got nil")
				return
			}
			dmp := diffmatchpatch.New()
			expected := `alias	index	routing.index	routing.search	is_write_index
foo	foo01	FooRoutingIndex	FooRoutingSearch	FooIsWriteIndex
bar	bar02	BarRoutingIndex	BarRoutingSearch	BarIsWriteIndex
`
			diffs := dmp.DiffMain(expected, *res, false)
			if dmp.DiffLevenshtein(diffs) > 0 {
				t.Errorf("detected diff:\n%s", dmp.DiffPrettyText(diffs))
			}
		},
	)
}
