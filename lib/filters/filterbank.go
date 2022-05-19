package filters

import (
	"fmt"
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
			name: "delta",
			noArgsFunMaker: func() func(s string) (string, error) {
				var last float64

				return func(s string) (string, error) {
					if s == "" {
						return "", nil
					}

					v, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return "", err
					}

					delta := v - last

					last = v

					formatted := fmt.Sprintf("%f", delta)

					return formatted, nil
				}
			},
		},
		&FilterEntry{
			name: "uniquecount",
			makerI: func(fixedWindowSize int) (func(string) (string, error), error) {
				circularbuffer := make([]string, fixedWindowSize)
				seen := map[string]int{}
				filling := true
				index := 0

				return func(value string) (string, error) {
					if filling {
						seen[value] += 1
						circularbuffer[index] = value
						index++

						if index >= fixedWindowSize {
							filling = false
							index = 0
						}
					} else {
						overwritten := circularbuffer[index]
						if overwritten != value {
							if seen[overwritten] == 1 {
								delete(seen, overwritten)
							} else {
								seen[overwritten] -= 1
							}

							circularbuffer[index] = value
							seen[value] += 1
						}
						index = (index + 1) % fixedWindowSize
					}

					currentUniq := len(seen)
					return strconv.Itoa(currentUniq), nil
				}, nil
			},
		},
		&FilterEntry{
			name: "rollingmean",
			makerI: func(fixedWindowSize int) (func(string) (string, error), error) {
				circularbuffer := make([]float64, fixedWindowSize)
				var total float64
				var count int
				index := 0

				return func(value string) (string, error) {
					var fvalue float64

					if value != "" {
						v, err := strconv.ParseFloat(value, 64)
						if err != nil {
							return "", err
						}
						fvalue = v
					}

					if count < fixedWindowSize {
						circularbuffer[index] = fvalue
						count++
						index++

						if index >= fixedWindowSize {
							index = 0
						}
					} else {
						overwritten := circularbuffer[index]
						circularbuffer[index] = fvalue
						total += fvalue - overwritten
						index = (index + 1) % fixedWindowSize
					}

					currentMean := total / float64(count)
					formatted := fmt.Sprintf("%f", currentMean)
					return formatted, nil
				}, nil
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
