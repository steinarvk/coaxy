package accessor

import "github.com/steinarvk/coaxy/lib/interfaces"

type Func func(interfaces.Record) (interfaces.Record, error)

func (h Func) Extract(rec interfaces.Record) (interfaces.Record, error) {
	return h(rec)
}
