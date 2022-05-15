package coaxy

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/google/shlex"
	"github.com/steinarvk/coaxy/lib/sniff"
)

type tupleReader struct {
	readTuple func() ([]string, error)
}

func (t *tupleReader) ReadTuple() ([]string, error) {
	return t.readTuple()
}

func makeCSVishReader(separator rune, hasHeader bool, r io.Reader) (*tupleReader, error) {
	reader := csv.NewReader(r)
	reader.Comma = separator
	reader.ReuseRecord = true

	if hasHeader {
		_, err := reader.Read()
		if err != nil {
			return nil, err
		}
	}

	rv := &tupleReader{
		readTuple: reader.Read,
	}

	return rv, nil
}

func makeShlexReader(hasHeader bool, r io.Reader) (*tupleReader, error) {
	scanner := bufio.NewScanner(r)

	readline := func() ([]string, error) {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			items, err := shlex.Split(line)
			if err != nil {
				return nil, err
			}

			return items, nil
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		return nil, io.EOF
	}

	if hasHeader {
		readline()
	}

	rv := &tupleReader{
		readTuple: readline,
	}

	return rv, nil
}

func makeTupleReader(desc *sniff.Descriptor, r io.Reader) (*tupleReader, error) {
	switch {
	case desc.Format == sniff.FormatCSV:
		return makeCSVishReader(',', desc.HasHeader, r)
	case desc.Format == sniff.FormatTSV:
		return makeCSVishReader('\t', desc.HasHeader, r)
	case desc.Format == sniff.FormatSSV:
		return makeCSVishReader(' ', desc.HasHeader, r)
	case desc.Format == sniff.FormatShell:
		return makeShlexReader(desc.HasHeader, r)
	default:
		return nil, fmt.Errorf("unsupported format %q for tuple-based reader", desc.Format)
	}
}
