package coaxy

import (
	"fmt"
	"io"
	"strconv"

	"github.com/steinarvk/coaxy/lib/accessor"
	"github.com/steinarvk/coaxy/lib/coaxyexpr"
	"github.com/steinarvk/coaxy/lib/interfaces"
	"github.com/steinarvk/coaxy/lib/prereader"
	"github.com/steinarvk/coaxy/lib/record"
	"github.com/steinarvk/coaxy/lib/sniff"
)

var (
	peekSize = 100 * 1024
)

func (s *Stream) Descriptor() *sniff.Descriptor {
	return s.descriptor
}

func OpenStream(r io.Reader) (*Stream, error) {
	header, r, err := prereader.Preread(r, peekSize)
	if err != nil {
		return nil, err
	}

	descriptor, err := sniff.Sniff(header)
	if err != nil {
		return nil, err
	}

	rv := &Stream{
		descriptor: descriptor,
	}

	if descriptor.TupleBased {
		tupleReader, err := makeTupleReader(descriptor, r)
		if err != nil {
			return nil, err
		}

		var indexByName map[string]int

		if len(descriptor.ColumnNames) > 0 {
			indexByName = map[string]int{}

			for i, name := range descriptor.ColumnNames {
				indexByName[name] = i
			}
		}

		readrecord := func() (interfaces.Record, error) {
			tuple, err := tupleReader.ReadTuple()
			if err != nil {
				return nil, err
			}

			return record.FromStringTuple(tuple, indexByName)
		}

		rv.readRecord = readrecord
	} else {
		switch descriptor.Format {
		case sniff.FormatJSON:
			rv.readRecord = makeJSONArrayReader(r)

		case sniff.FormatJSONL:
			rv.readRecord = makeJSONLRecordReader(r)

		default:
			return nil, fmt.Errorf("unsupported non-tuple-based format: %q", descriptor.Format)
		}
	}

	return rv, nil
}

func (s *Stream) ResolveField(query string) (interfaces.Accessor, error) {
	rv, err := s.resolveField(query)
	if err != nil {
		return nil, fmt.Errorf("invalid field accessor %q: %w", query, err)
	}
	return rv, nil
}

func (s *Stream) resolveField(query string) (interfaces.Accessor, error) {
	desc := s.descriptor

	n, err := strconv.Atoi(query)
	if err == nil && n >= 0 {
		if !desc.TupleBased {
			return nil, fmt.Errorf("numeric column reference for non-tuple-based format %q", desc.Format)
		}

		if n >= desc.NumColumns {
			return nil, fmt.Errorf("out of bounds for %q with %d columns", desc.Format, desc.NumColumns)
		}

		return accessor.AtIndex(n), nil
	}

	if desc.TupleBased {
		for index, name := range desc.ColumnNames {
			if name == query {
				return accessor.AtIndex(index), nil
			}
		}
	}

	expr, err := coaxyexpr.Parse(query)
	if err != nil {
		return nil, err
	}

	return expr.MakeAccessor(), nil
}

func (s *Stream) Select(columns []interfaces.Accessor) (*Reader, error) {
	if s.selected {
		return nil, fmt.Errorf("Select() was already called")
	}
	s.selected = true

	readrow := func(out []string) error {
		if len(out) != len(columns) {
			return fmt.Errorf("reading tuple of length %d; got buffer of size %d", len(columns), len(out))
		}

		rec, err := s.readRecord()
		if err != nil {
			return err
		}

		for outindex, accessor := range columns {
			valuerec, err := accessor.Extract(rec)
			if err != nil {
				return err
			}

			value, err := valuerec.AsValue()
			if err != nil {
				return err
			}

			out[outindex] = value
		}

		return nil
	}

	return &Reader{
		tupleSize: len(columns),
		read:      readrow,
	}, nil
}

func (r *Reader) Read(out []string) error {
	return r.read(out)
}

func (r *Reader) ForEach(f func([]string) error) error {
	for {
		tuple := make([]string, r.tupleSize)

		if err := r.Read(tuple); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if err := f(tuple); err != nil {
			return err
		}
	}
}
