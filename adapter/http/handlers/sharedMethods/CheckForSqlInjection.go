package sharedMethods

import (
	"regexp"

	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/sharedMethods"
)

type CheckForSqlInjection struct {
	customError port.CustomErrorInterface
}

func NewCheckForSqlInjection(
	customError port.CustomErrorInterface,
) sharedMethods.CheckForSqlInjectionInterface {
	return &CheckForSqlInjection{
		customError: customError,
	}
}

func (checkForSqlInjection *CheckForSqlInjection) CheckForSqlInjection(input string) error {
	sqlInjectionPattern := `(?i)\b(?:SELECT\s[^*]*\*\sFROM|INSERT\sINTO|UPDATE\s|DELETE\sFROM|TRUNCATE\sTABLE|DROP\sTABLE|ALTER\sTABLE|UNION\sALL|UNION\sSELECT|EXEC(?:UTE){0,1}|MERGE|CALL|CAST\()|(\b(?:SELECT|INSERT|UPDATE|DELETE|DROP|UNION|TRUNCATE|EXEC|DECLARE|DATABASE|ALTER|CREATE|XP_)\b)|(--\s|/\*|\*/)`

	match, err := regexp.MatchString(sqlInjectionPattern, input)
	if err != nil {
		return checkForSqlInjection.customError.ThrowError(err.Error())
	}

	if match {
		return checkForSqlInjection.customError.ThrowError(
			helper.ErrorWithCodeConst,
			helper.ErrorWhenTryToExecuteRowsQueryCodeConst,
		)
	}

	return nil
}
