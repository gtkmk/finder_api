package userRequestEntity

import (
	"encoding/json"
	"net/http"

	"github.com/gtkmk/finder_api/infra/EmptyBodyVerification"
)

type EditUserRequest struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	CellphoneNumber   string `json:"cellphone_number"`
	PermissionGroupId string `json:"permission_group_id"`
	IsActive          bool   `json:"is_active"`
}

func DecodeEditUserRequest(req *http.Request) (*EditUserRequest, error) {
	if err := EmptyBodyVerification.ValidateBody(req); err != nil {
		return nil, err
	}

	var user *EditUserRequest

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (user *EditUserRequest) Validate() error {
	idValidationError := ValidateId(user.Id)

	if idValidationError != nil {
		return idValidationError
	}

	nameValidationError := ValidateName(user.Name)

	if nameValidationError != nil {
		return nameValidationError
	}

	emailValidationError := ValidateEmail(user.Email)

	if emailValidationError != nil {
		return emailValidationError
	}

	if user.CellphoneNumber != "" {
		cellphoneNumberValidationError := ValidateCellphoneNumber(user.CellphoneNumber)

		if cellphoneNumberValidationError != nil {
			return cellphoneNumberValidationError
		}
	}

	return nil
}
