package main

import (
	"testing"
)

func testDebugFunc(msg string) {}

func TestLoadConfigFile(t *testing.T) {
	t.Run(
		"file note exist",
		func(t *testing.T) {
			var cfg *EsConfig
			err := cfg.LoadFromFile(".", "notfound", testDebugFunc)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if cfg != nil {
				t.Errorf("expected nil, got %v", *cfg)
			}
		},
	)
}
