package sniff

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func makeCSVishParser(formatType FormatType, separator rune) func(r io.Reader) (*Descriptor, error) {
	return func(r io.Reader) (*Descriptor, error) {
		reader := csv.NewReader(r)
		reader.Comma = separator
		reader.ReuseRecord = false

		var records [][]string

		for {
			rec, err := reader.Read()
			if err == io.EOF {
				break
			}
			records = append(records, rec)
		}

		if len(records) <= 1 {
			return nil, nil
		}

		commonLength := len(records[0])
		for _, rec := range records {
			if len(rec) != commonLength {
				return nil, nil
			}
		}

		if commonLength <= 1 {
			return nil, nil
		}

		columnrecords := make([][]string, commonLength)
		for _, rec := range records {
			for i, value := range rec {
				columnrecords[i] = append(columnrecords[i], value)
			}
		}

		hasHeader := false

		var typesIfNoHeader []ValueType
		var typesIfHeader []ValueType

		for _, values := range columnrecords {
			typeIfNoHeader := DetectValueType(values)
			typeIfHeader := DetectValueType(values[1:])

			typesIfNoHeader = append(typesIfNoHeader, typeIfNoHeader)
			typesIfHeader = append(typesIfHeader, typeIfHeader)

			if typeIfNoHeader.Kind != typeIfHeader.Kind && typeIfNoHeader.Kind == KindString {
				hasHeader = true
			}
		}

		desc := &Descriptor{
			Format:     formatType,
			HasHeader:  hasHeader,
			NumColumns: commonLength,
			FieldTypes: map[string]ValueType{},
		}

		var fieldNames []string
		var types []ValueType

		if hasHeader {
			desc.ColumnNames = records[0]
			fieldNames = records[0]
			types = typesIfHeader
		} else {
			types = typesIfNoHeader
			for i, _ := range types {
				fieldNames = append(fieldNames, fmt.Sprintf("%d", i+1))
			}
		}

		for i, name := range fieldNames {
			desc.FieldTypes[name] = types[i]
		}

		return desc, nil
	}
}

func sniffLines(lines []string) (*Descriptor, error) {
	joined := strings.Join(lines, "")

	detectors := []func(io.Reader) (*Descriptor, error){
		makeCSVishParser(FormatCSV, ','),
		makeCSVishParser(FormatTSV, '\t'),
		makeCSVishParser(FormatSSV, ' '),
	}

	for _, detector := range detectors {
		desc, err := detector(bytes.NewBufferString(joined))
		if err != nil {
			return nil, err
		}
		if desc != nil {
			return desc, nil
		}
	}

	return nil, nil
}
