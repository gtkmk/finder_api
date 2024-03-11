package credentialsDomain

import (
	"strings"
)

type Credential struct {
	Email    string
	Password string
}

func NewCredentials(email, password string) *Credential {
	email = strings.ToLower(email)

	return &Credential{
		email,
		password,
	}
}
