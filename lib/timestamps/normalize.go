package timestamps

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type parser struct {
	re     *regexp.Regexp
	layout string
	custom func(string) (time.Time, bool, error)
}

func (p *parser) tryParse(s string) (time.Time, bool, error) {
	if p.custom != nil {
		return p.custom(s)
	}

	if p.re != nil {
		if !p.re.MatchString(s) {
			return time.Time{}, false, nil
		}
	}

	t, err := time.Parse(p.layout, s)
	if err != nil {
		if p.re == nil {
			err = nil
		}
		return time.Time{}, false, err
	}

	return t, true, nil
}

type multiparser struct {
	normalizers []func(string) string
	parsers     []*parser
}

func (m *multiparser) tryParse(s string) (time.Time, bool, error) {
	if s == "" {
		return time.Time{}, false, nil
	}

	original := s

	for j, norm := range m.normalizers {
		s := norm(original)

		if j > 0 && s == original {
			continue
		}

		for i, p := range m.parsers {
			t, ok, err := p.tryParse(s)
			if err != nil {
				return time.Time{}, false, err
			}

			if !ok {
				continue
			}

			if i != 0 {
				tmp := m.parsers[0]
				m.parsers[0] = p
				m.parsers[i] = tmp
			}

			if j != 0 {
				tmp := m.normalizers[0]
				m.normalizers[0] = norm
				m.normalizers[j] = tmp
			}

			return t, true, nil
		}
	}

	return time.Time{}, false, nil
}

func (m *multiparser) parse(s string) (time.Time, error) {
	t, ok, err := m.tryParse(s)
	if err != nil {
		return time.Time{}, err
	}

	if !ok {
		return time.Time{}, fmt.Errorf("failed to parse %q as timestamp", s)
	}

	return t, nil
}

var (
	minTimestamp = time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	maxTimestamp = time.Date(2050, 1, 1, 0, 0, 0, 0, time.UTC)
)

func makeIntParser(unit time.Duration) func(string) (time.Time, bool, error) {
	toNanos := int64(unit) / int64(time.Nanosecond)
	t0 := int64(minTimestamp.UnixNano()) / int64(unit)
	t1 := int64(maxTimestamp.UnixNano()) / int64(unit)

	return func(s string) (time.Time, bool, error) {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return time.Time{}, false, nil
		}

		if n < t0 || n > t1 {
			return time.Time{}, false, nil
		}

		n *= toNanos

		t := time.Unix(0, n).UTC()
		return t, true, nil
	}
}

func makeNumberParser(unit time.Duration) func(string) (time.Time, bool, error) {
	toNanos := float64(unit) / float64(time.Nanosecond)
	t0 := float64(minTimestamp.UnixNano()) / float64(unit)
	t1 := float64(maxTimestamp.UnixNano()) / float64(unit)

	return func(s string) (time.Time, bool, error) {
		x, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return time.Time{}, false, nil
		}

		if x < t0 || x > t1 {
			return time.Time{}, false, nil
		}

		x *= toNanos

		t := time.Unix(0, int64(x)).UTC()
		return t, true, nil
	}
}

func nullNormalizer(s string) string {
	return s
}

func colonToCommaNormalizer(s string) string {
	if strings.Count(s, ":") < 3 {
		return s
	}

	i := strings.LastIndex(s, ":")
	return s[:i] + "." + s[i+1:]
}

var (
	standardNormalizers = []func(string) string{
		nullNormalizer,
		colonToCommaNormalizer,
	}

	standardFormats = []*parser{
		&parser{
			custom: makeIntParser(time.Second),
		},
		&parser{
			custom: makeIntParser(time.Millisecond),
		},
		&parser{
			custom: makeIntParser(time.Microsecond),
		},
		&parser{
			custom: makeIntParser(time.Nanosecond),
		},
		&parser{
			custom: makeNumberParser(time.Second),
		},
		&parser{
			custom: makeNumberParser(time.Millisecond),
		},
		&parser{
			custom: makeNumberParser(time.Microsecond),
		},
		&parser{
			custom: makeNumberParser(time.Nanosecond),
		},
		&parser{
			layout: time.RFC3339Nano,
		},
		&parser{
			layout: time.RFC3339,
		},
		&parser{
			layout: time.ANSIC,
		},
		&parser{
			layout: time.RFC822Z,
		},
		&parser{
			layout: time.RFC822,
		},
		&parser{
			layout: time.UnixDate,
		},
		&parser{
			layout: time.RubyDate,
		},
		&parser{
			layout: time.RFC850,
		},
		&parser{
			layout: time.RFC1123Z,
		},
		&parser{
			layout: time.RFC1123,
		},
		&parser{
			layout: "2006-01-02 15:04:05Z07:00",
		},
		&parser{
			layout: "Jan 2 2006 03:04:05PM",
		},
		&parser{
			layout: "Jan 2 2006 15:04:05",
		},
	}
)

func newMultiparser() *multiparser {
	normalizers := make([]func(string) string, len(standardNormalizers))
	for i, p := range standardNormalizers {
		normalizers[i] = p
	}

	parsers := make([]*parser, len(standardFormats))
	for i, p := range standardFormats {
		parsers[i] = p
	}

	return &multiparser{
		normalizers: normalizers,
		parsers:     parsers,
	}
}

func newNormalizerFunc(f func(time.Time) string) func(string) (string, error) {
	p := newMultiparser()

	return func(s string) (string, error) {
		t, err := p.parse(s)
		if err != nil {
			return "", err
		}

		return f(t), nil
	}
}

func NewNormalizerUnix() func(string) (string, error) {
	return newNormalizerFunc(func(t time.Time) string {
		unixSecs := float64(t.UnixNano()) / 1e9
		return fmt.Sprintf("%f", unixSecs)
	})
}

func NewNormalizerISO() func(string) (string, error) {
	return newNormalizerFunc(func(t time.Time) string {
		return t.Format(time.RFC3339Nano)
	})
}
