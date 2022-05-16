package filters

import (
	"strconv"

	"github.com/steinarvk/coaxy/lib/timestamps"
)

func makeBank(entries ...*FilterEntry) map[string]*FilterEntry {
	m := map[string]*FilterEntry{}
	for _, entry := range entries {
		m[entry.name] = entry
	}
	return m
}

var (
	filterbank = makeBank(
		&FilterEntry{
			name: "isnull",
			noArgsFun: func(s string) (string, error) {
				if s == "" {
					return "true", nil
				}
				return "false", nil
			},
		},
		&FilterEntry{
			name: "categorize",
			noArgsFunMaker: func() func(s string) (string, error) {
				m := map[string]int{}

				return func(s string) (string, error) {
					n, ok := m[s]
					if !ok {
						n = len(m) + 1
						m[s] = n
					}

					return strconv.Itoa(n), nil
				}
			},
		},
		&FilterEntry{
			name:           "unixtime",
			noArgsFunMaker: timestamps.NewNormalizerUnix,
		},
		&FilterEntry{
			name:           "isotime",
			noArgsFunMaker: timestamps.NewNormalizerISO,
		},
	)
)
