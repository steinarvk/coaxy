package record

import (
	"github.com/steinarvk/chaxy/lib/sniff"
)

type Stream struct {
	descriptor *sniff.Descriptor
	readRecord func() (record, error)
	selected   bool
}

type Reader struct {
	tupleSize int
	read      func([]string) error
}
