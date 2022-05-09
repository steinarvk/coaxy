package accessor

import "github.com/steinarvk/coaxy/lib/interfaces"

type Error struct {
	Err error
}

func (e Error) Extract(rec interfaces.Record) (interfaces.Record, error) {
	return nil, e.Err
}
