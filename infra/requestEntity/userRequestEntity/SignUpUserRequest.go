package userRequestEntity

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gtkmk/finder_api/core/domain/cpf"
	"github.com/gtkmk/finder_api/core/domain/phone"
	"github.com/gtkmk/finder_api/infra/EmptyBodyVerification"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

type SignUpUserRequest struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	UserName        string `json:"user_name"`
	Email           string `json:"email"`
	Cpf             string `json:"cpf"`
	Password        string `json:"password"`
	CellphoneNumber string `json:"cellphone_number"`
	IsActive        bool   `json:"is_active"`
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
	MaximumPasswordLengthConst  = 50
	UserFieldNameConst          = "O nome do usuário"
	UserFieldUserNameConst      = "O nome de usuário"
	UserFieldCpfConst           = "O cpf do usuário"
	UserFieldPasswordConst      = "A senha do usuário"
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
	fmt.Print("******************************************************")
	fmt.Print(signUpUserRequest)
	cpfValidationError := ValidateCpf(signUpUserRequest.Cpf)

	if cpfValidationError != nil {
		return cpfValidationError
	}

	nameValidationError := ValidateName(signUpUserRequest.Name)

	if nameValidationError != nil {
		return nameValidationError
	}

	userNameValidationError := ValidateUserName(signUpUserRequest.UserName)

	if userNameValidationError != nil {
		return userNameValidationError
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

	passwordValidationError := ValidatePassword(signUpUserRequest.Password)

	if passwordValidationError != nil {
		return passwordValidationError
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

func ValidateUserName(userName string) error {
	nameValidationError := requestEntityFieldsValidation.ValidateField(
		userName,
		UserFieldUserNameConst,
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

func ValidatePassword(password string) error {
	nameValidationError := requestEntityFieldsValidation.ValidatePasswordField(
		password,
		UserFieldPasswordConst,
		MaximumPasswordLengthConst,
	)

	if nameValidationError != nil {
		return nameValidationError
	}

	return nil
}
