package record

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/steinarvk/chaxy/lib/prereader"
	"github.com/steinarvk/chaxy/lib/sniff"
)

const (
	peekSize = 4096
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

		rv.tupleReader = tupleReader
	} else {
		return nil, fmt.Errorf("not implemented: non-tuple-based formats")
	}

	return rv, nil
}

func (s *Stream) ResolveField(query string) (*Accessor, error) {
	rv, err := s.resolveField(query)
	if err != nil {
		return nil, fmt.Errorf("invalid field accessor %q: %w", query, err)
	}
	return rv, nil
}

func (s *Stream) resolveField(query string) (*Accessor, error) {
	desc := s.descriptor

	n, err := strconv.Atoi(query)
	if err == nil && n >= 0 {
		if !desc.TupleBased {
			return nil, fmt.Errorf("numeric column reference for non-tuple-based format %q", desc.Format)
		}

		if n == 0 {
			return nil, errors.New("column references are 1-based")
		}

		if n > desc.NumColumns {
			return nil, fmt.Errorf("out of bounds for %q with %d columns", desc.Format, desc.NumColumns)
		}

		return &Accessor{
			byIndex: true,
			index:   n - 1,
		}, nil
	}

	if desc.TupleBased && len(desc.ColumnNames) == 0 {
		return nil, errors.New("tuple-based format without field names; must use numeric column reference")
	}

	for index, name := range desc.ColumnNames {
		if name == query {
			return &Accessor{
				byIndex: true,
				index:   index,
			}, nil
		}
	}

	// no more advanced queries supported currently.
	// to come: path-based, regex-based, as well as transformers (e.g. "categorized integer").

	return nil, errors.New("no such column")
}

func (s *Stream) Select(columns []*Accessor) (*Reader, error) {
	if s.selected {
		return nil, fmt.Errorf("Select() was already called")
	}
	s.selected = true

	if s.tupleReader == nil {
		return nil, fmt.Errorf("not implemented: non-tuple-based readers")
	}

	indices := make([]int, len(columns))
	for i, col := range columns {
		if !col.byIndex {
			return nil, fmt.Errorf("not implemented: non-index-based accessors for tuple-based formats")
		}
		indices[i] = col.index
	}

	readrow := func(out []string) error {
		if len(out) != len(indices) {
			return fmt.Errorf("reading tuple of length %d; got buffer of size %d", len(indices), len(out))
		}

		tuple, err := s.tupleReader.ReadTuple()
		if err != nil {
			return err
		}

		for outindex, inindex := range indices {
			out[outindex] = tuple[inindex]
		}

		return nil
	}

	return &Reader{
		tupleSize: len(indices),
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
