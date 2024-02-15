package port

import "github.com/gtkmk/finder_api/core/domain/email"

type EmailServiceInterface interface {
	SendEmail(emailStruct email.Email) error
}
