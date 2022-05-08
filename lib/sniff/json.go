package sniff

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"

	"github.com/bcicen/jstream"
	"github.com/steinarvk/chaxy/lib/chaxyvalue"
)

func jsonlParser(r io.Reader) (*Descriptor, error) {
	scanner := bufio.NewScanner(r)

	var objects []interface{}

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

		objects = append(objects, generic)
	}

	// Read error?
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return autodetectFromRecords(FormatJSONL, objects)
}

func autodetectFromRecords(format FormatType, objects []interface{}) (*Descriptor, error) {
	valuesSeen := map[string][]string{}

	for _, generic := range objects {
		flattened, err := flattenObject(generic)
		if err != nil {
			return nil, err
		}

		for key, value := range flattened {
			stringified, err := chaxyvalue.JSONPrimitiveToString(value)
			if err != nil {
				return nil, err
			}

			valuesSeen[key] = append(valuesSeen[key], stringified)
		}
	}

	if len(valuesSeen) == 0 {
		return nil, nil
	}

	desc := &Descriptor{
		Format:     format,
		FieldTypes: map[string]ValueType{},
	}

	for key, values := range valuesSeen {
		desc.FieldTypes[key] = DetectValueType(values)
	}

	return desc, nil
}

func jsonArrayParser(r io.Reader) (*Descriptor, error) {
	var objects []interface{}

	decoder := jstream.NewDecoder(r, 1)
	for rec := range decoder.Stream() {
		objects = append(objects, rec.Value)
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
