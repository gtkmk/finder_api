package phone

import (
	"github.com/gtkmk/finder_api/core/domain/helper"
	"regexp"
	"strconv"
	"strings"

	
)

const (
	PHONE_MIN_LENGTH = 10
	PHONE_MAX_LENGTH = 11
)

type PHONE struct {
	number string
}

func NewPhone(number string) *PHONE {
	return &PHONE{
		number,
	}
}

func (phone *PHONE) AsString() string {
	return phone.number
}

func (phone *PHONE) AsUnmaskedString() (string, error) {
	if err := phone.Validate(); err != nil {
		return "", err
	}

	return phone.number, nil
}

func (phone *PHONE) Validate() error {
	phone.Unmask()

	if len(phone.number) != 0 {
		if err := phone.isValid(); err != nil {
			return err
		}

		_, err := strconv.Atoi(phone.number)

		if err != nil {
			return helper.ErrorBuilder(helper.ThePhoneShouldOnlyCountNumbersConst)
		}
	}

	return nil
}

func (phone *PHONE) Unmask() {
	replacer := strings.NewReplacer(" ", "", "(", "", ")", "", "-", "", "+", "")
	phone.number = replacer.Replace(phone.number)
}

func (phone *PHONE) isValid() error {
	if err := phone.checkLength(); err != nil {
		return err
	}

	if err := phone.checkOnlyNumbers(); err != nil {
		return err
	}

	return nil
}

func (phone *PHONE) checkLength() error {
	if !phone.isCorrectLength() {
		return helper.ErrorBuilder(helper.ThePhoneIsNotInAvalidLengthConst)
	}

	return nil
}

func (phone *PHONE) isCorrectLength() bool {
	return len(phone.number) >= PHONE_MIN_LENGTH && len(phone.number) <= PHONE_MAX_LENGTH
}

func (phone *PHONE) checkOnlyNumbers() error {
	if !phone.hasOnlyNumbers() {
		return helper.ErrorBuilder(helper.ThePhoneShouldOnlyCountNumbersConst)
	}

	return nil
}

func (phone *PHONE) hasOnlyNumbers() bool {
	regex := regexp.MustCompile(`^\d+$`)

	return regex.MatchString(phone.number)
}

func (phone *PHONE) AddMask() string {
	re := regexp.MustCompile(`(\d{2})(\d{4,5})(\d{4})`)
	return re.ReplaceAllString(phone.number, "($1) $2-$3")
}
