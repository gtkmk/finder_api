package userRequestEntity

import (
	"encoding/json"
	"net/http"

	"github.com/gtkmk/finder_api/infra/EmptyBodyVerification"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

type SignInUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignInDecodeUserRequest(req *http.Request) (*SignInUserRequest, error) {
	if err := EmptyBodyVerification.ValidateBody(req); err != nil {
		return nil, err
	}

	var user *SignInUserRequest
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		return nil, err
	}

	return user, requestEntityFieldsValidation.ValidateEmailField(user.Email)
}
