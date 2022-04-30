package main

import (
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestBuildTemplateParametersWithStreamError(t *testing.T) {
	tempFile, err := ioutil.TempFile(os.TempDir(), "cles-template-temp-")
	if err != nil {
		t.Error(err)
	}
	defer tempFile.Close()

	_, err = tempFile.WriteString("INVALID JSON")
	if err != nil {
		t.Error(err)
	}
	tempFile.Seek(0, io.SeekStart)

	got, err := buildTemplateParameters("foo", "", tempFile)
	if err == nil {
		t.Errorf("expected error, got %v", got)
	}
}

func TestBuildTemplateParametersWithEOFError(t *testing.T) {
	tempFile, err := ioutil.TempFile(os.TempDir(), "cles-template-temp-")
	if err != nil {
		t.Error(err)
	}
	tempFile.Close()

	got, err := buildTemplateParameters("foo", "", tempFile)
	if err == nil {
		t.Errorf("expected error, got %v", got)
	}
}

func TestBuildTemplateParametersWithStream(t *testing.T) {
	tempFile, err := ioutil.TempFile(os.TempDir(), "cles-template-temp-")
	if err != nil {
		t.Error(err)
	}
	defer tempFile.Close()

	_, err = tempFile.WriteString("{\"q\":\"SEARCH\"}")
	if err != nil {
		t.Error(err)
	}
	tempFile.Seek(0, io.SeekStart)

	actual, err := buildTemplateParameters("foo-template", "", tempFile)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expected := map[string]interface{}{
		"id": "foo-template",
		"params": map[string]interface{}{
			"q": "SEARCH",
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestBuildTemplateParametersWithRawParams(t *testing.T) {
	actual, err := buildTemplateParameters("foo-template", "{\"q\":\"SEARCH\"}", os.Stdin)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expected := map[string]interface{}{
		"id": "foo-template",
		"params": map[string]interface{}{
			"q": "SEARCH",
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
