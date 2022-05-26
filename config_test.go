package main

import (
	"testing"
)

func testDebugFunc(msg string) {}

func TestLoadConfigFile(t *testing.T) {
	t.Run(
		"file note exist",
		func(t *testing.T) {
			got, err := loadConfigFile(".", "notfound", testDebugFunc)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if got != nil {
				t.Errorf("expected nil, got %v", got)
			}
		},
	)
}
