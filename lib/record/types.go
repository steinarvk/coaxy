package record

import (
	"github.com/steinarvk/chaxy/lib/sniff"
)

type Stream struct {
	descriptor  *sniff.Descriptor
	tupleReader *tupleReader
	selected    bool
}

type Accessor struct {
	byIndex bool
	index   int
}

type Reader struct {
	read func([]string) error
}
