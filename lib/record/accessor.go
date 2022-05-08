package record

import (
	"errors"

	"github.com/steinarvk/coaxy/lib/interfaces"
)

type simpleAccessor struct {
	byIndex      bool
	bySimpleName bool
	index        int
	name         string
}

func (a *simpleAccessor) Extract(rec interfaces.Record) (interfaces.Record, error) {
	switch {
	case a.byIndex:
		return rec.GetByIndex(a.index)

	case a.bySimpleName:
		return rec.GetByName(a.name)

	default:
		return nil, errors.New("invalid accessor")
	}
}
