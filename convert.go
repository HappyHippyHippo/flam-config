package config

import (
	"fmt"
	"strings"

	flam "github.com/happyhippyhippo/flam"
)

func Convert(
	val any,
) any {
	if pValue, ok := val.(flam.Bag); ok {
		result := flam.Bag{}
		for k, value := range pValue {
			result[strings.ToLower(k)] = Convert(value)
		}

		return result
	}

	if lValue, ok := val.([]any); ok {
		var result []any
		for _, i := range lValue {
			result = append(result, Convert(i))
		}

		return result
	}

	if mValue, ok := val.(map[string]any); ok {
		result := flam.Bag{}
		for k, i := range mValue {
			result[strings.ToLower(k)] = Convert(i)
		}

		return result
	}

	if mValue, ok := val.(map[any]any); ok {
		result := flam.Bag{}
		for k, i := range mValue {
			stringKey, ok := k.(string)
			if ok {
				result[strings.ToLower(stringKey)] = Convert(i)
			} else {
				result[fmt.Sprintf("%v", k)] = Convert(i)
			}
		}

		return result
	}

	if fValue, ok := val.(float64); ok && float64(int(fValue)) == fValue {
		return int(fValue)
	}

	return val
}
