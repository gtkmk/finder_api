package requestEntityFieldsValidation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gtkmk/finder_api/core/domain/helper"

	"github.com/google/uuid"
	emailDomain "github.com/gtkmk/finder_api/core/domain/email"
)

const (
	notANumberConst = "NaN"
)

var (
	uppercaseRegex   = regexp.MustCompile(`[A-Z]`)
	lowercaseRegex   = regexp.MustCompile(`[a-z]`)
	numberRegex      = regexp.MustCompile(`[0-9]`)
	specialCharRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

func ValidateField(field, fieldName string, maxLength int) error {
	if len(field) > maxLength || len(field) == 0 {
		if len(field) > maxLength {
			return helper.ErrorBuilder(helper.FieldCannotHaveMoreThanSetCharactersConst, fieldName, maxLength)
		}

		return helper.ErrorBuilder(helper.FieldCannotBeEmptyConst, fieldName)
	}

	return nil
}

func ValidateDateField(field, fieldName string, maxLength int, layoutDateRegex string) error {
	if len(field) > maxLength || len(field) == 0 {
		switch {
		case strings.Contains(field, notANumberConst):
			return helper.ErrorBuilder(helper.FieldCannotBeEmptyConst, fieldName)
		case len(field) > maxLength:
			return helper.ErrorBuilder(helper.FieldCannotHaveMoreThanSetCharactersConst, fieldName, maxLength)
		default:
			return helper.ErrorBuilder(helper.FieldCannotBeEmptyConst, fieldName)
		}
	}

	_, err := regexp.MatchString(layoutDateRegex, field)

	if err != nil {
		return helper.ErrorBuilder(helper.InvalidDateFormatConst)
	}

	return nil
}

func ValidateFieldInArray(field, fieldName string, allowedValues []string) error {
	if len(field) == 0 {
		return helper.ErrorBuilder(helper.FieldCannotBeEmptyConst, fieldName)
	}

	for _, val := range allowedValues {
		fmt.Println(val, " ==:> ", field)
		if field == val {
			return nil
		}
	}

	return helper.ErrorBuilder(helper.FieldNotInAllowedValuesConst, fieldName)
}

func IsValidUUID(fieldname string, hash string) error {
	_, err := uuid.Parse(hash)

	if err != nil {
		return helper.ErrorBuilder(helper.InformFieldConst, fieldname)
	}

	return nil
}

func ValidateEmailField(email string) error {
	emailValidator := emailDomain.NewEmail(email)

	return emailValidator.Validate()
}

func ValidateGroupLayerField(layer int64) error {
	if layer == 0 {
		return helper.ErrorBuilder(helper.TheGroupPermissionLayerCannotBeResetConst)
	}

	if layer < 0 {
		return helper.ErrorBuilder(helper.TheGroupPermissionLayerCannotBeNegativeConst)
	}

	if layer > 999 {
		return helper.ErrorBuilder(helper.TheGroupPermissionLayerCannotBeGreaterThanLimitConst)
	}

	return nil
}

func ValidatePasswordField(field, fieldName string, maxLength int) error {
	if len(field) > maxLength || len(field) == 0 {
		if len(field) > maxLength {
			return helper.ErrorBuilder(helper.FieldCannotHaveMoreThanSetCharactersConst)
		}
		return helper.ErrorBuilder(helper.FieldCannotBeEmptyConst)
	}

	if !uppercaseRegex.MatchString(field) {
		return helper.ErrorBuilder(helper.PasswordMissingUppercaseConst)
	}

	if !lowercaseRegex.MatchString(field) {
		return helper.ErrorBuilder(helper.PasswordMissingLowercaseConst)
	}

	if !numberRegex.MatchString(field) {
		return helper.ErrorBuilder(helper.PasswordMissingNumberConst)
	}

	if !specialCharRegex.MatchString(field) {
		return helper.ErrorBuilder(helper.PasswordMissingSpecialCharConst)
	}

	return nil
}
