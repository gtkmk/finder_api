package passwordDomain

import (
	"github.com/gtkmk/finder_api/core/domain/helper"

	"regexp"
)

type Password struct {
	value string
}

const PasswordMinimumLengthConst = 8

func NewPassWord(value string) *Password {
	return &Password{
		value,
	}
}

func (password *Password) Validate() error {
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	if !lowercaseRegex.MatchString(password.value) {

		return helper.ErrorBuilder(helper.PasswordMustHaveCharactersMessageConst)
	}

	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	if !uppercaseRegex.MatchString(password.value) {
		return helper.ErrorBuilder(helper.PasswordMustHaveAtLastOneUpperCaseMessageConst)
	}

	digitRegex := regexp.MustCompile(`\d`)
	if !digitRegex.MatchString(password.value) {
		return helper.ErrorBuilder(helper.PasswordMustHaveAtLastOneDigitMessageConst)
	}

	specialCharRegex := regexp.MustCompile(`[@$!%*?&]`)
	if !specialCharRegex.MatchString(password.value) {
		return helper.ErrorBuilder(helper.PasswordMustHaveAtLastOneSpecialCharacterConst)
	}

	if len(password.value) < PasswordMinimumLengthConst {
		return helper.ErrorBuilder(helper.PasswordMustHaveTheMinLengthMessageConst, PasswordMinimumLengthConst)
	}

	return nil
}
