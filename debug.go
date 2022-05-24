package main

import (
	"os"

	"github.com/fatih/color"
)

func getDebugStream(debug bool) (*os.File, error) {
	var stream *os.File
	if debug {
		stream = os.Stdout
		return stream, nil
	}
	stream, err := os.OpenFile(os.DevNull, os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func initDebugFunc(debug bool) (*func(message string), error) {
	stream, err := getDebugStream(debug)
	if err != nil {
		return nil, err
	}
	f := color.New(color.FgCyan).FprintFunc()
	debugFunc := func(message string) {
		f(stream, message+"\n")
	}
	return &debugFunc, nil
}
