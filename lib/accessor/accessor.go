package accessor

import (
	"github.com/steinarvk/coaxy/lib/interfaces"
)

type AtIndex int

func (n AtIndex) Extract(rec interfaces.Record) (interfaces.Record, error) {
	return rec.GetByIndex(int(n))
}

type AtField string

func (k AtField) Extract(rec interfaces.Record) (interfaces.Record, error) {
	return rec.GetByName(string(k))
}

type Chain []interfaces.Accessor

func (c Chain) Extract(rec interfaces.Record) (interfaces.Record, error) {
	for _, acc := range c {
		newrec, err := acc.Extract(rec)
		if err != nil {
			return nil, err
		}
		rec = newrec
	}

	return rec, nil
}
