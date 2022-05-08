package chaxyvalue

import "fmt"

const (
	maxSafeJSONInteger = 9007199254740991
	minSafeJSONInteger = -9007199254740991
)

func JSONPrimitiveToString(value interface{}) (string, error) {
	switch v := value.(type) {
	case int:
		return fmt.Sprintf("%d", v), nil

	case string:
		return v, nil

	case bool:
		if v {
			return "true", nil
		}
		return "false", nil

	case float64:
		asInteger := int64(v)
		if minSafeJSONInteger <= asInteger && asInteger <= maxSafeJSONInteger {
			if float64(asInteger) == v {
				return fmt.Sprintf("%d", asInteger), nil
			}
		}

		return fmt.Sprintf("%v", v), nil

	case nil:
		// For the purposes of this tool -- compatible with CSV etc. -- we use the empty string as our null value.
		return "", nil

	default:
		return "", fmt.Errorf("not a JSON primitive: %v", value)
	}
}
