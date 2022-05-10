package sniff

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"

	"github.com/steinarvk/coaxy/lib/interfaces"
	"github.com/steinarvk/coaxy/lib/record"
)

func makeCSVishParser(formatType FormatType, separator rune) func(r io.Reader) (*Descriptor, error) {
	return func(r io.Reader) (*Descriptor, error) {
		reader := csv.NewReader(r)
		reader.Comma = separator
		reader.ReuseRecord = false

		var records []interfaces.Record
		var tuples [][]string

		var commonLength int

		for {
			rec, err := reader.Read()
			if err == io.EOF {
				break
			}

			if commonLength == 0 {
				commonLength = len(rec)
			} else {
				if commonLength != len(rec) {
					return nil, nil
				}
			}

			truerec, err := record.FromStringTuple(rec, nil)
			if err != nil {
				return nil, err
			}

			tuples = append(tuples, rec)
			records = append(records, truerec)
		}

		if len(records) <= 1 {
			return nil, nil
		}

		columnrecords := make([][]string, commonLength)
		for _, rec := range tuples {
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
			TupleBased: true,
			HasHeader:  hasHeader,
			NumColumns: commonLength,
		}

		if hasHeader {
			desc.ColumnNames = tuples[0]

			var newrecs []interfaces.Record

			names := map[string]int{}
			for i, name := range desc.ColumnNames {
				names[name] = i
			}

			for _, rec := range records[1:] {
				newrec, err := record.StringTupleWithNames(rec, names)
				if err != nil {
					return nil, err
				}
				newrecs = append(newrecs, newrec)
			}

			records = newrecs
		}

		nestedTypes, err := DetectNestedFields(records)
		if err != nil {
			return nil, err
		}

		desc.FieldTypes = nestedTypes

		return desc, nil
	}
}

func sniffLines(lines []string) (*Descriptor, error) {
	joined := strings.Join(lines, "")

	detectors := []func(io.Reader) (*Descriptor, error){
		makeCSVishParser(FormatCSV, ','),
		makeCSVishParser(FormatTSV, '\t'),
		makeCSVishParser(FormatSSV, ' '),
		jsonlParser,
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
