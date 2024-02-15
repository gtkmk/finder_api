package repositories

import (
	"github.com/gtkmk/finder_api/core/domain/userDomain"
)

type UserRepository interface {
	VerifyIfUserExistsByCpf(cpf string) bool
	FindUserByEmail(email string) (*userDomain.User, error)
	FindUserById(id string) (*userDomain.User, error)
	CreateUser(user *userDomain.User) error
	UpdateResetPasswordStatus(toggle bool, status string, userId string) error
}
