package sniff

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"

	"github.com/bcicen/jstream"
	"github.com/steinarvk/coaxy/lib/interfaces"
	"github.com/steinarvk/coaxy/lib/record"
)

func jsonlParser(r io.Reader) (*Descriptor, error) {
	scanner := bufio.NewScanner(r)

	var records []interfaces.Record

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		var generic interface{}

		if err := json.Unmarshal([]byte(line), &generic); err != nil {
			// Not JSONLines
			return nil, nil
		}

		rec, err := record.FromJSONValue(generic)
		if err != nil {
			return nil, err
		}

		records = append(records, rec)
	}

	// Read error?
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return autodetectFromRecords(FormatJSONL, records)
}

func autodetectFromRecords(format FormatType, objects []interfaces.Record) (*Descriptor, error) {
	nestedTypes, err := DetectNestedFields(objects)
	if err != nil {
		return nil, err
	}

	if len(nestedTypes) == 0 {
		return nil, nil
	}

	desc := &Descriptor{
		Format:     format,
		FieldTypes: nestedTypes,
	}

	return desc, nil
}

func jsonArrayParser(r io.Reader) (*Descriptor, error) {
	var objects []interfaces.Record

	decoder := jstream.NewDecoder(r, 1)
	for val := range decoder.Stream() {
		rec, err := record.FromJSONValue(val.Value)
		if err != nil {
			return nil, err
		}

		objects = append(objects, rec)
	}

	if len(objects) < 2 {
		return nil, nil
	}

	// last one might be partial
	objects = objects[:len(objects)-1]

	// ignore errors!
	_ = decoder.Err()

	return autodetectFromRecords(FormatJSON, objects)
}
