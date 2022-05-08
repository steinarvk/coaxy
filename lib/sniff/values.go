package sniff

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	arbitraryStringValue = ValueType{
		Kind: KindString,
	}
)

func isNull(s string) bool {
	return s == ""
}

func parseAsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func parseAsNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func parseAsBool(s string) bool {
	s = strings.ToLower(s)
	return s == "true" || s == "false" || s == "t" || s == "f"
}

func allSatisfy(xs []string, f func(string) bool) bool {
	for _, x := range xs {
		if !f(x) {
			return false
		}
	}

	return true
}

func someSatisfy(xs []string, f func(string) bool) bool {
	for _, x := range xs {
		if f(x) {
			return true
		}
	}

	return false
}

func filterValuesNot(xs []string, f func(string) bool) []string {
	var ys []string

	for _, x := range xs {
		if !f(x) {
			ys = append(ys, x)
		}
	}

	return ys
}

type pattern struct {
	Pattern *regexp.Regexp
	Type    ValueType
}

var (
	patterns = []pattern{
		pattern{
			Pattern: regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`),
			Type: ValueType{
				Kind:   KindDate,
				Format: "YYYY-MM-DD",
			},
		},
		pattern{
			Pattern: regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}[T ][0-9]{2}:[0-9]{2}(:[0-9]{2}([.][0-9]+)?)?(Z|([0-9]{2}:[0-9]{2}))?$`),
			Type: ValueType{
				Kind:   KindTimestamp,
				Format: "ISO8601",
			},
		},
	}
)

func DetectValueType(values []string) ValueType {
	if len(values) <= 0 {
		return arbitraryStringValue
	}

	var retval ValueType

	if someSatisfy(values, isNull) {
		retval.Optional = true
		values = filterValuesNot(values, isNull)

		if len(values) == 0 {
			return ValueType{Kind: KindNull}
		}
	}

	switch {
	case allSatisfy(values, parseAsInt):
		retval.Kind = KindInt

	case allSatisfy(values, parseAsNumber):
		retval.Kind = KindNumber

	default:
		invalidated := make([]bool, len(patterns))

		for _, value := range values {
			for i, pattern := range patterns {
				if invalidated[i] {
					continue
				}

				if !pattern.Pattern.MatchString(value) {
					invalidated[i] = true
				}
			}
		}

		for i, pattern := range patterns {
			if !invalidated[i] {
				result := pattern.Type
				result.Optional = retval.Optional
				return result
			}
		}

		retval.Kind = KindString
	}

	return retval
}
