package userRequestEntity

import (
	"encoding/json"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"net/http"

	"github.com/gtkmk/finder_api/core/domain/passwordDomain"
	
	"github.com/gtkmk/finder_api/infra/EmptyBodyVerification"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

type ResetUserPassword struct {
	Id              string `json:"id"`
	OldPassword     string `json:"old_password"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func ResetUserPasswordRequest(req *http.Request) (*ResetUserPassword, error) {
	if err := EmptyBodyVerification.ValidateBody(req); err != nil {
		return nil, err
	}

	var request *ResetUserPassword
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func (resetUserPassword *ResetUserPassword) Validate() error {
	if err := requestEntityFieldsValidation.IsValidUUID(UserIdFieldConst, resetUserPassword.Id); err != nil {
		return err
	}

	if resetUserPassword.OldPassword == "" {
		return helper.ErrorBuilder(helper.IsNecessaryToSendTheOldPasswordMessageConst)
	}

	if resetUserPassword.Password == "" || resetUserPassword.ConfirmPassword == "" {
		return helper.ErrorBuilder(helper.PasswordMayNotEmptyMessageConst)
	}

	if resetUserPassword.Password != resetUserPassword.ConfirmPassword {
		return helper.ErrorBuilder(helper.PasswordsNotEqualMessageConst)
	}

	passwordStrengthValidator := passwordDomain.NewPassWord(resetUserPassword.Password)

	return passwordStrengthValidator.Validate()
}
