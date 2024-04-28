package helper

func RetrieveSafeInterfaceValue(value interface{}, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}
	return value
}
