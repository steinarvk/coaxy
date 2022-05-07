package record

import (
	"io"

	"github.com/steinarvk/chaxy/lib/sniff"
)

type Stream struct {
	descriptor       *sniff.Descriptor
	reader           io.Reader
	shouldSkipHeader bool
}

type Accessor struct {
}

type Record struct {
}
