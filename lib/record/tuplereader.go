package record

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/steinarvk/chaxy/lib/sniff"
)

type tupleReader struct {
	reader *csv.Reader
}

func (t *tupleReader) ReadTuple() ([]string, error) {
	return t.reader.Read()
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
		reader: reader,
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
	default:
		return nil, fmt.Errorf("unsupported format %q for tuple-based reader", desc.Format)
	}
}
