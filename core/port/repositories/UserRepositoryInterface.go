package repositories

import (
	"github.com/gtkmk/finder_api/core/domain/userDomain"
)

type UserRepository interface {
	VerifyIfUserExistsByCpf(cpf string) bool
	VerifyIfUserExistsByUserName(userName string) bool
	FindUserByEmail(email string) (*userDomain.User, error)
	FindUserById(id string) (*userDomain.User, error)
	CreateUser(user *userDomain.User) error
	UpdateResetPasswordStatus(toggle bool, status string, userId string) error
	SetUserStatus(userId string, status string) error
	ResetUserPassword(userId string, password string) error
	FindCompleteUserInfoByID(userId string) ([]map[string]interface{}, error)
}
