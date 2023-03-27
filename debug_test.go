package main

import (
	"os"
	"testing"
)

func TestGetDebugStream(t *testing.T) {
	t.Run(
		"enable debug",
		func(t *testing.T) {
			stream, err := getDebugStream(true)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if stream != os.Stdout {
				t.Errorf("expected os.Stdout, got %v", stream)
			}
		},
	)
	t.Run(
		"disable debug",
		func(t *testing.T) {
			stream, err := getDebugStream(false)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if stream == os.Stdout {
				t.Errorf("expected os.DevNull, got os.Stdout")
			}
		},
	)
}

func TestInitDebugFunc(t *testing.T) {
	t.Run(
		"test debug func",
		func(t *testing.T) {
			fp, err := initDebugFunc(true)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			fn := *fp
			fn("debug")
		},
	)
}
