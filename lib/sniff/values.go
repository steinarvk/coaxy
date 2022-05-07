package sniff

import "strconv"

var (
	arbitraryStringValue = ValueType{
		Kind: KindString,
	}
)

func parseAsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func parseAsNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func allSatisfy(xs []string, f func(string) bool) bool {
	for _, x := range xs {
		if !f(x) {
			return false
		}
	}

	return true
}

func DetectValueType(values []string) ValueType {
	if len(values) <= 0 {
		return arbitraryStringValue
	}

	switch {
	case len(values) <= 0:
		return arbitraryStringValue

	case allSatisfy(values, parseAsInt):
		return ValueType{
			Kind: KindInt,
		}

	case allSatisfy(values, parseAsNumber):
		return ValueType{
			Kind: KindNumber,
		}

	default:
		return arbitraryStringValue
	}
}
