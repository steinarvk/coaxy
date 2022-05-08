package sniff

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"

	"github.com/steinarvk/chaxy/lib/chaxyvalue"
)

func jsonlParser(r io.Reader) (*Descriptor, error) {
	scanner := bufio.NewScanner(r)

	valuesSeen := map[string][]string{}

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

	// Read error?
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(valuesSeen) == 0 {
		return nil, nil
	}

	desc := &Descriptor{
		Format:     FormatJSONL,
		FieldTypes: map[string]ValueType{},
	}

	for key, values := range valuesSeen {
		desc.FieldTypes[key] = DetectValueType(values)
	}

	return desc, nil
}
