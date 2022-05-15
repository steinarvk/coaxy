package sniff

import (
	"bufio"
	"io"
	"strings"

	"github.com/google/shlex"
	"github.com/steinarvk/coaxy/lib/interfaces"
	"github.com/steinarvk/coaxy/lib/record"
)

func shlexParser(r io.Reader) (*Descriptor, error) {
	scanner := bufio.NewScanner(r)

	var records []interfaces.Record
	var tuples [][]string

	var commonLength int

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		rec, err := shlex.Split(line)
		if err != nil {
			return nil, nil
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

	// Read error?
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(records) <= 1 {
		return nil, nil
	}

	if commonLength <= 1 {
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
		Format:     FormatShell,
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
