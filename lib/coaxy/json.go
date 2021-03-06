package coaxy

import (
	"io"

	"github.com/bcicen/jstream"
	"github.com/steinarvk/coaxy/lib/interfaces"
	"github.com/steinarvk/coaxy/lib/record"
)

func makeJSONArrayReader(r io.Reader) func() (interfaces.Record, error) {
	decoder := jstream.NewDecoder(r, 1)

	stream := decoder.Stream()

	return func() (interfaces.Record, error) {
		rec, ok := <-stream

		if err := decoder.Err(); err != nil {
			return nil, err
		}

		if !ok {
			return nil, io.EOF
		}

		return record.FromJSONValue(rec.Value)
	}
}
