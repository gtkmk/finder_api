package userRequestEntity

import (
	"encoding/json"
	"net/http"

	"github.com/gtkmk/finder_api/core/domain/cpf"
	"github.com/gtkmk/finder_api/core/domain/phone"
	"github.com/gtkmk/finder_api/infra/EmptyBodyVerification"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

type SignUpUserRequest struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	Cpf               string `json:"cpf"`
	CellphoneNumber   string `json:"cellphone_number"`
	CreatorId         string `json:"creator_id"`
	PermissionGroupId string `json:"permission_group_id"`
	IsActive          bool   `json:"is_active"`
}

const (
	UserIdConst = "um usuário válido"
)

const (
	MaximumNameLengthConst      = 100
	MaximumCpfLengthConst       = 11
	MaximumBirthDateLengthConst = 10
	MaximumCellPhoneLengthConst = 12
	MaximumEmailLengthConst     = 50
	UserFieldNameConst          = "O nome do usuário"
	UserFieldCpfConst           = "O cpf do usuário"
	UserFieldBirthDateConst     = "A data de nascimento do usuário"
	UserFieldEmailConst         = "O email do usuário"
	UserFieldCellphoneConst     = "O número do celular do usuário"
	UserFieldPhoneConst         = "O número do telefone do usuário"
	UserFieldUserIdConst        = "o usuário"
)

func SignUpDecodeUserRequest(req *http.Request) (*SignUpUserRequest, error) {
	if err := EmptyBodyVerification.ValidateBody(req); err != nil {
		return nil, err
	}

	var user *SignUpUserRequest

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (signUpUserRequest *SignUpUserRequest) Validate() error {
	cpfValidationError := ValidateCpf(signUpUserRequest.Cpf)

	if cpfValidationError != nil {
		return cpfValidationError
	}

	nameValidationError := ValidateName(signUpUserRequest.Name)

	if nameValidationError != nil {
		return nameValidationError
	}

	emailValidationError := ValidateEmail(signUpUserRequest.Email)

	if emailValidationError != nil {
		return emailValidationError
	}

	if signUpUserRequest.CellphoneNumber != "" {
		cellphoneNumberValidationError := ValidateCellphoneNumber(signUpUserRequest.CellphoneNumber)

		if cellphoneNumberValidationError != nil {
			return cellphoneNumberValidationError
		}
	}

	return nil
}

func ValidateId(id string) error {
	idValidationError := requestEntityFieldsValidation.IsValidUUID(
		UserIdConst,
		id,
	)

	if idValidationError != nil {
		return idValidationError
	}

	return nil
}

func ValidateCellphoneNumber(cellphoneNumber string) error {
	cellphone := phone.NewPhone(cellphoneNumber)

	if cellphoneValidationError := cellphone.Validate(); cellphoneValidationError != nil {
		return cellphoneValidationError
	}

	return nil
}

func ValidateCpf(userCPf string) error {
	cpf := cpf.NewCpf(userCPf)

	if cpfValidationError := cpf.Validate(); cpfValidationError != nil {
		return cpfValidationError
	}

	return nil
}

func ValidateName(name string) error {
	nameValidationError := requestEntityFieldsValidation.ValidateField(
		name,
		UserFieldNameConst,
		MaximumNameLengthConst,
	)

	if nameValidationError != nil {
		return nameValidationError
	}

	return nil
}

func ValidateEmail(email string) error {
	nameValidationError := requestEntityFieldsValidation.
		ValidateEmailField(email)

	if nameValidationError != nil {
		return nameValidationError
	}

	return nil
}
