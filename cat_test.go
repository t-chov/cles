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

type catIndicesImpl struct {
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

func (c catIndicesImpl) Do(ctx context.Context) (elastic.CatIndicesResponse, error) {
	if c.ExpectError {
		return []elastic.CatIndicesResponseRow{}, fmt.Errorf("cat indices error")
	}
	resp := []elastic.CatIndicesResponseRow{
		{
			Health:       "green",
			Status:       "true",
			Index:        "foo",
			UUID:         "XXXX",
			Pri:          1,
			Rep:          0,
			DocsCount:    1024,
			DocsDeleted:  8,
			StoreSize:    "64KB",
			PriStoreSize: "64.8KB",
		},
		{
			Health:       "yellow",
			Status:       "true",
			Index:        "bar",
			UUID:         "YYYY",
			Pri:          2,
			Rep:          3,
			DocsCount:    2048,
			DocsDeleted:  128,
			StoreSize:    "1MB",
			PriStoreSize: "8MB",
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

func TestPrettyCatIndices(t *testing.T) {
	t.Run(
		"fail catIndices API",
		func(t *testing.T) {
			svc := catIndicesImpl{
				ExpectError: true,
			}
			res, err := prettyCatIndices(svc)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
			if !strings.Contains(fmt.Sprintf("%s", err), "cat indices error") {
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
			svc := catIndicesImpl{
				ExpectError: false,
			}
			res, err := prettyCatIndices(svc)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if res == nil {
				t.Errorf("expected not nil, got nil")
				return
			}
			dmp := diffmatchpatch.New()
			expected := `health	status	index	uuid	pri	rep	docs.count	docs.deleted	store.size	pri.store.size
green	true	foo	XXXX	1	0	1024	8	64KB	64.8KB
yellow	true	bar	YYYY	2	3	2048	128	1MB	8MB
`
			diffs := dmp.DiffMain(expected, *res, false)
			if dmp.DiffLevenshtein(diffs) > 0 {
				t.Errorf("detected diff:\n%s", dmp.DiffPrettyText(diffs))
			}
		},
	)
}
