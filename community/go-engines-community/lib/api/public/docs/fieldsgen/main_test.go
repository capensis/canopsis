package main

import (
	"bytes"
	"os"
	"testing"
)

func TestGenerate_GivenCommittedFile_ShouldBeUpToDate(t *testing.T) {
	expected, err := os.ReadFile("../fields_swagger.yaml")
	if err != nil {
		t.Fatal(err)
	}

	got, err := generate()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, got) {
		t.Error("fields_swagger.yaml is outdated, run go generate ./cmd/public-api/")
	}
}
