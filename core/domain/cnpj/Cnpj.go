package cnpj

import (
	"github.com/gtkmk/finder_api/core/domain/helper"
	"regexp"
	"strconv"
	"strings"

	
)

const (
	CnpjLengthConst                = 14
	FirstPositionVerificationConst = 11
)

type CNPJ struct {
	Number string
}

func NewCnpj(number string) *CNPJ {
	return &CNPJ{
		number,
	}
}

func (cnpj *CNPJ) AsString() string {
	return cnpj.Number
}

func (cnpj *CNPJ) Validate() error {
	cnpj.Number = regexp.MustCompile(`\D`).ReplaceAllString(cnpj.Number, "")

	if err := cnpj.isValid(); err != nil {
		return err
	}

	_, err := strconv.Atoi(cnpj.Number)

	if err != nil {
		return helper.ErrorBuilder(helper.CnpjMustHaveOnlyNumberMessageConst)
	}

	return nil
}

func IsCnpj(number string) error {
	return NewCnpj(number).Validate()
}

func (cnpj *CNPJ) isValid() error {
	if err := cnpj.checkLength(); err != nil {
		return err
	}

	if err := cnpj.checkOnlyNumbers(); err != nil {
		return err
	}

	if err := cnpj.zeroSequences(); err != nil {
		return err
	}

	if err := cnpj.checkVerificationDigits(); err != nil {
		return err
	}

	return nil
}

func (cnpj *CNPJ) checkLength() error {
	if !cnpj.isCorrectLength() {
		return helper.ErrorBuilder(helper.InvalidCnpjLengthMessageConst)
	}

	return nil
}

func (cnpj *CNPJ) isCorrectLength() bool {
	return len(cnpj.Number) == CnpjLengthConst
}

func (cnpj *CNPJ) checkOnlyNumbers() error {
	if !cnpj.hasOnlyNumbers() {
		return helper.ErrorBuilder(helper.SendAValidCNPJErrorMessageConst)
	}

	return nil
}

func (cnpj *CNPJ) hasOnlyNumbers() bool {
	regex := regexp.MustCompile(`^\d+$`)

	return regex.MatchString(cnpj.Number)
}

func (cnpj *CNPJ) zeroSequences() error {
	if cnpj.Number == "00000000000000" {
		return helper.ErrorBuilder(helper.InvalidCnpjWithOnlyZerosMessageConst)
	}

	return nil
}

func (cnpj *CNPJ) checkVerificationDigits() error {
	for i := 0; i <= 1; i++ {
		var numero int

		j := 5 + i
		soma := 0

		for numero = 0; numero <= (FirstPositionVerificationConst + i); numero++ {
			if numero == (4 + i) {
				j = 9
			}

			result := []rune(cnpj.Number)
			character := int(result[numero]) - '0'
			soma += character * j
			j--
		}

		resto := soma % 11
		var digito int

		if resto < 2 {
			digito = 0
		} else {
			digito = 11 - resto
		}

		result := []rune(cnpj.Number)
		character := string(result[12+i])

		if character != strconv.Itoa(digito) {
			return helper.ErrorBuilder(helper.SendAValidCNPJErrorMessageConst)
		}
	}
	return nil
}

func (cnpj *CNPJ) Sanitize() {
	cnpj.Number = strings.Replace(cnpj.Number, ".", "", -1)
	cnpj.Number = strings.Replace(cnpj.Number, "-", "", -1)
	cnpj.Number = strings.Replace(cnpj.Number, "/", "", -1)
}
