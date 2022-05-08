package record

import (
	"github.com/steinarvk/coaxy/lib/interfaces"
	"github.com/steinarvk/coaxy/lib/sniff"
)

type Stream struct {
	descriptor *sniff.Descriptor
	readRecord func() (interfaces.Record, error)
	selected   bool
}

type Reader struct {
	tupleSize int
	read      func([]string) error
}
