package util

import "strconv"

func ConvertAnyToString(v interface{}) string {
	switch result := v.(type) {
	case bool:
		return strconv.FormatBool(result)
	case float64:
		return strconv.FormatFloat(result, 'f', -1, 64)
	case string:
		return result
	default:
		return ""
	}
}
