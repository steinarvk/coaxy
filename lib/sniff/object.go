package sniff

import (
	"errors"
	"fmt"
	"strings"
)

func flattenObject(root interface{}) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	var err error

	var visit func(key string, obj interface{})
	visit = func(key string, obj interface{}) {
		switch value := obj.(type) {
		case bool:
			m[key] = obj
		case float64:
			m[key] = obj
		case string:
			m[key] = obj
		case nil:
			m[key] = obj

		case []interface{}:
			for i, element := range value {
				subkey := fmt.Sprintf("%s[%d]", key, i)
				visit(subkey, element)
			}

		case map[string]interface{}:
			for k, v := range value {
				subkey := strings.TrimLeft(key+"."+k, ".")
				visit(subkey, v)
			}

		default:
			err = errors.New("invalid JSON value")
		}
	}

	visit("", root)

	if err != nil {
		return nil, err
	}

	return m, nil
}
