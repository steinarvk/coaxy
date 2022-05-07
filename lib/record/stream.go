package record

import (
	"io"

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

	return &Stream{
		descriptor:       descriptor,
		reader:           r,
		shouldSkipHeader: descriptor.HasHeader,
	}, nil
}
