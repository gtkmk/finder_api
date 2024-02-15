package requestEntityFieldsValidation

import (
	"github.com/gtkmk/finder_api/core/domain/helper"
	"regexp"
	"strings"

	"github.com/google/uuid"
	emailDomain "github.com/gtkmk/finder_api/core/domain/email"
	
)

const (
	notANumberConst = "NaN"
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
