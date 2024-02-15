package helper

import (
	"fmt"
	"strconv"
)

func ConvertToFloat(value interface{}) (*float64, error) {
	defaultResult := 0.0

	switch result := value.(type) {
	case float64:
		return &result, nil
	case string:
		return getFloatFromString(result)
	case nil:
		return &defaultResult, nil
	default:
		return nil, fmt.Errorf(ErrorConvertingValueConst, "float64")
	}
}

func getFloatFromString(value string) (*float64, error) {
	floatValue, err := strconv.ParseFloat(value, 64)

	if err != nil {

		return nil, fmt.Errorf(ErrorConvertingValueConst, "float64")
	}

	return &floatValue, nil
}
