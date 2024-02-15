package userRequestEntity

import (
	"encoding/json"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"net/http"

	"github.com/gtkmk/finder_api/core/domain/passwordDomain"
	
	"github.com/gtkmk/finder_api/infra/EmptyBodyVerification"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

type ForgotUserPassword struct {
	Id              string `json:"id"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func ForgotUserPasswordRequest(req *http.Request) (*ForgotUserPassword, error) {
	if err := EmptyBodyVerification.ValidateBody(req); err != nil {
		return nil, err
	}

	var request *ForgotUserPassword
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func (forgotUserPassword *ForgotUserPassword) ThrowsErrorIfPasswordsDoesNotMatch() error {
	if err := requestEntityFieldsValidation.IsValidUUID(UserIdFieldConst, forgotUserPassword.Id); err != nil {
		return err
	}

	if forgotUserPassword.Password == "" || forgotUserPassword.ConfirmPassword == "" {
		return helper.ErrorBuilder(helper.PasswordMayNotEmptyMessageConst)
	}

	if forgotUserPassword.Password != forgotUserPassword.ConfirmPassword {
		return helper.ErrorBuilder(helper.PasswordsAreDifferentsConst)
	}

	passwordStrengthValidator := passwordDomain.NewPassWord(forgotUserPassword.Password)

	return passwordStrengthValidator.Validate()
}
