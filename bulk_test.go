package main

import (
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestIdToStrWithError(t *testing.T) {
	id := 1.234
	got, err := idToStr(id)
	if err == nil {
		t.Errorf("isToStr(%v) must fail. actual: %v", id, got)
	}
}

func TestIdToStr(t *testing.T) {
	tests := []struct {
		name     string
		arg      interface{}
		expected string
	}{
		{
			name:     "int convert",
			arg:      11,
			expected: "11",
		},
		{
			name:     "string convert",
			arg:      "foobar",
			expected: "foobar",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := idToStr(test.arg)
			if err != nil {
				t.Errorf("idToStr(%v) must success. error: %v", test.arg, err)
			}
			if got != test.expected {
				t.Errorf("idToStr(%v) expected %v. got: %v", test.arg, test.expected, got)
			}
		})
	}
}

func TestCalsNumOfLines(t *testing.T) {
	tempFile, err := ioutil.TempFile(os.TempDir(), "cles-test-temp-")
	if err != nil {
		t.Error(err)
	}
	defer tempFile.Close()

	_, err = tempFile.WriteString("foo\nbar\nbaz")
	if err != nil {
		t.Error(err)
	}

	tempFile.Seek(0, io.SeekStart)
	actual := calcNumOfLines(tempFile)

	if actual != 3 {
		t.Errorf("expected 3 lines. got: %d", actual)
	}
}

func TestRowToDoc(t *testing.T) {
	docString := "{\"one\":1, \"two\":\"bar\", \"three\":false, \"four\":true, \"five\":[0, 1, 2]}"
	actual, err := rowToDoc([]byte(docString))
	if err != nil {
		t.Error(err)
	}
	expected := map[string]interface{}{
		"one":   1,
		"two":   "bar",
		"three": false,
		"four":  true,
		"five":  []int{0, 1, 2},
	}
	if reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestRowToDocError(t *testing.T) {
	docString := "INVALID JSON"
	actual, err := rowToDoc([]byte(docString))
	if err == nil {
		t.Errorf("expected error. got %v", actual)
	}
}

func TestGetId(t *testing.T) {
	doc := map[string]interface{}{
		"foo_id": 100,
		"bar_id": "200",
	}
	tests := []struct {
		name     string
		idColumn string
		expected string
	}{
		{
			name:     "input int ID",
			idColumn: "foo_id",
			expected: "100",
		},
		{
			name:     "input string ID",
			idColumn: "bar_id",
			expected: "200",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := getId(doc, test.idColumn)
			if err != nil {
				t.Errorf("expected no error. error: %v", err)
			}
			if *got != test.expected {
				t.Errorf("expected %v. got: %v", test.expected, *got)
			}
		})
	}
}

func TestGetIdNil(t *testing.T) {
	doc := map[string]interface{}{
		"foo_id": 100,
		"bar_id": "200",
	}
	got, err := getId(doc, "")
	if err != nil {
		t.Errorf("expected no error. error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil. got: %v", err)
	}
}

func TestGetIdError(t *testing.T) {
	doc := map[string]interface{}{
		"foo_id": 100,
		"bar_id": "200",
		"baz_id": []int{0, 1, 2},
	}
	tests := []struct {
		name     string
		idColumn string
	}{
		{
			name:     "invalid datatype",
			idColumn: "baz_id",
		},
		{
			name:     "not exist key",
			idColumn: "qux_id",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := getId(doc, test.idColumn)
			if err == nil {
				t.Errorf("expected error. got: %v", got)
			}
		})
	}
}
