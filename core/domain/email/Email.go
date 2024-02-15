package email

import (
	"github.com/gtkmk/finder_api/core/domain/helper"
	"regexp"
	"strings"

	
)

const (
	EmailMaxLengthConst = 50
)

type Email struct {
	Email string
}

func NewEmail(email string) *Email {
	email = strings.ToLower(email)

	return &Email{
		email,
	}
}

func (email *Email) AsString() string {
	return email.Email
}

func (email *Email) Validate() error {
	if len(email.Email) == 0 {
		return helper.ErrorBuilder(helper.EmailCannotBeEmptyConst)
	}

	if len(email.Email) > EmailMaxLengthConst {
		return helper.ErrorBuilder(helper.EmailTooLongConst)
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email.Email) {
		return helper.ErrorBuilder(helper.InvalidEmailFormatConst)
	}

	return nil
}

func IsValidEmail(number string) bool {
	if err := NewEmail(number).Validate(); err != nil {
		return false
	}

	return true
}
