package sniff

import (
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

const (
	testFileDir = "../../testdata"
)

func TestDetectTypeOfFiles(t *testing.T) {
	infos, err := ioutil.ReadDir(testFileDir)
	if err != nil {
		t.Fatal(err)
	}

	wantTestTypes := []string{
		"csv",
		"tsv",
		"ssv",
		"json",
		"jsonl",
	}

	testedTypes := map[string]int{}

	for _, info := range infos {
		var expectedType string

		switch {
		case strings.HasSuffix(info.Name(), ".csv"):
			expectedType = "csv"

		case strings.HasSuffix(info.Name(), ".tsv"):
			expectedType = "tsv"

		case strings.HasSuffix(info.Name(), ".ssv"):
			expectedType = "ssv"

		case strings.HasSuffix(info.Name(), ".jsonl"):
			expectedType = "jsonl"

		case strings.HasSuffix(info.Name(), ".json"):
			expectedType = "json"

		default:
			continue
		}

		testedTypes[expectedType] += 1

		fullPath := path.Join(testFileDir, info.Name())
		data, err := ioutil.ReadFile(fullPath)
		if err != nil {
			t.Fatalf("error reading %q: %v", fullPath, err)
		}

		desc, err := Sniff(data)
		if err != nil {
			t.Errorf("error sniffing %q: %v", fullPath, err)
			continue
		}

		if desc == nil {
			t.Errorf("error sniffing %q: returned nil", fullPath)
		}

		got := string(desc.Format)
		want := expectedType

		if got != want {
			t.Errorf("error sniffing %q: got %q want %q", fullPath, got, want)
		}
	}

	for _, want := range wantTestTypes {
		if testedTypes[want] == 0 {
			t.Errorf("expected at least one test of type %q", want)
		}
	}
}
