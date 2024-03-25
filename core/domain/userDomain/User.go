package userDomain

import (
	"time"

	emailDomain "github.com/gtkmk/finder_api/core/domain/email"
)

const (
	UserStatusPendingConst = "pending"
	UserStatusLoggedConst  = "logged"
	UserStatusExpiredConst = "expired"
	PasswordResetConst     = "password_reset"
	FirstAccessConst       = "first_access"
	ProposalRejectedConst  = "proposal_rejected"
)

const (
	TranslatedFieldNameUserName = "Nome"
	TranslatedFieldNameUnity    = "Divisão"
	TranslatedFieldNameIsActive = "Ativo/Inativo"
	TranslatedFieldNameUserRole = "Perfil/Cargo"
	TranslatedIsActive          = "Ativo"
	TranslatedIsInactive        = "Inativo"
)

type User struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	UserName        string    `json:"userName"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	Cpf             string    `json:"cpf"`
	CellphoneNumber string    `json:"cellphone_number"`
	Status          string    `json:"status"`
	IsActive        bool      `json:"is_active"`
	ResetPassword   bool      `json:"reset_password"`
	CreatedAt       time.Time `json:"created_at"`
}

func NewUser(
	id string,
	name string,
	userName string,
	email string,
	password string,
	cpf string,
	cellphoneNumber string,
	status string,
	isActive bool,
	resetPassword bool,
	createdAt time.Time,
) *User {
	email = emailDomain.NewEmail(email).AsString()

	return &User{
		Id:              id,
		Name:            name,
		UserName:        userName,
		Email:           email,
		Password:        password,
		Cpf:             cpf,
		CellphoneNumber: cellphoneNumber,
		Status:          status,
		IsActive:        isActive,
		ResetPassword:   resetPassword,
		CreatedAt:       createdAt,
	}
}
