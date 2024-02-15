package userRequestEntity

import (
	"encoding/json"
	"net/http"

	"github.com/gtkmk/finder_api/infra/EmptyBodyVerification"
	"github.com/gtkmk/finder_api/infra/requestEntity/proposalRequestEntity"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

const (
	EmailFieldConst = "O Email do usuário"
)

func ForgotPasswordDecodeRequest(req *http.Request) (*ForgotPasswordRequest, error) {
	if err := EmptyBodyVerification.ValidateBody(req); err != nil {
		return nil, err
	}

	var user *ForgotPasswordRequest
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		return nil, err
	}

	return user, requestEntityFieldsValidation.ValidateField(
		user.Email,
		EmailFieldConst,
		proposalRequestEntity.MaximumEmailLengthConst,
	)
}
