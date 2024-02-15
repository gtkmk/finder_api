package credentialsDomain

import emailDomain "github.com/gtkmk/finder_api/core/domain/email"

type Credential struct {
	Email    string
	Password string
}

func NewCredentials(email, password string) *Credential {
	email = emailDomain.NewEmail(email).AsString()

	return &Credential{
		email,
		password,
	}
}
