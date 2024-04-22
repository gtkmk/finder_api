package helper

import (
	"fmt"
	"strconv"
)

func ConvertToInt64(value interface{}) (*int64, error) {
	defaultResult := int64(0)

	switch result := value.(type) {
	case int64:
		return &result, nil
	case string:
		return getInt64FromString(result)
	case nil:
		return &defaultResult, nil
	default:
		return nil, fmt.Errorf(ErrorConvertingValueConst, "int64")
	}
}

func getInt64FromString(value string) (*int64, error) {
	intValue, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		return nil, fmt.Errorf(ErrorConvertingValueConst, "int64")
	}

	return &intValue, nil
}
