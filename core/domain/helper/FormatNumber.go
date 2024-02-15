package helper

import (
	"fmt"
	"strings"
)

func FormatCurrency(value int) string {
	valueStr := fmt.Sprintf("%d", value)

	n := len(valueStr)

	if n <= 3 {
		return valueStr + ",00"
	}

	if n <= 6 {
		return valueStr[:n-3] + "." + valueStr[n-3:] + ",00"
	}

	return valueStr[:n-6] + "." + valueStr[n-6:n-3] + "." + valueStr[n-3:] + ",00"
}

func FormatCurrencyWithThousandsSeparator(value float64) string {
	if value == 0 {
		return "0"
	}

	formatted := fmt.Sprintf("%.2f", value)
	parts := strings.Split(formatted, ".")

	integerPart := parts[0]
	decimalPart := parts[1]
	count := 0

	var result string

	for i := len(integerPart) - 1; i >= 0; i-- {
		result = string(integerPart[i]) + result
		count++
		if count%3 == 0 && i != 0 {
			result = "." + result
		}
	}

	formattedValue := result + "," + decimalPart

	return formattedValue
}
