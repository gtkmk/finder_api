package helper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func ConvertToString(value interface{}) string {
	switch result := value.(type) {
	case string:
		return result
	case int64:
		return strconv.FormatInt(result, 10)
	case float64:
		return strconv.FormatFloat(result, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(result), 'f', -1, 32)
	case int:
		return strconv.FormatInt(int64(result), 10)
	case []byte:
		return string(result)
	case time.Time:
		return result.Format("2006-01-02 15:04:05")
	case map[string]interface{}:
		return getStringFromMapStringInterface(result)
	case nil:
		return ""
	default:
		fmt.Print(ErrorBuilder(
			ErrorConvertingValueConst,
			"string",
		))
		return ""
	}
}

func getStringFromMapStringInterface(value map[string]interface{}) string {
	jsonBytes, err := json.Marshal(value)

	if err != nil {
		fmt.Print(ErrorBuilder(
			ErrorConvertingValueConst,
			"string",
		))
		return ""
	}

	return string(jsonBytes)
}
