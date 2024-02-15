package repository

import (
	"database/sql"
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/userDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	"github.com/gtkmk/finder_api/infra/database/models"
	"strings"
)

type UserDatabase struct {
	connection port.ConnectionInterface
}

func NewUserDatabase(connection port.ConnectionInterface) repositories.UserRepository {
	return &UserDatabase{
		connection,
	}
}

func (userDatabase *UserDatabase) VerifyIfUserExistsByCpf(cpf string) bool {
	var exists bool
	query := `SELECT CASE WHEN EXISTS (SELECT 1 FROM user WHERE cpf = ?) THEN 1 ELSE 0 END`

	if err := userDatabase.connection.Raw(query, &exists, cpf); err != nil {
		return false
	}

	return exists
}

func (userDatabase *UserDatabase) FindUserByEmail(email string) (*userDomain.User, error) {
	var dbUser *models.UserWithPermissionGroup

	query := `SELECT
        id as user_id,
        name as user_name,
        email as user_email,
		password as user_password,
		cpf as user_cpf,
        role as user_role,
		cellphone_number as user_cellphone_number,
        is_active as user_is_active,
        status as user_status,
        created_at as user_created_at,
      FROM user
      WHERE email = ?`

	if err := userDatabase.connection.Raw(query, &dbUser, email); err != nil {
		return nil, err
	}

	if dbUser == nil {
		return nil, nil
	}

	user := userDomain.NewUser(
		dbUser.UserID,
		dbUser.UserName,
		dbUser.UserEmail,
		dbUser.UserPassword,
		dbUser.UserCpf,
		dbUser.UserCellphoneNumber,
		dbUser.UserStatus,
		dbUser.UserIsActive,
		dbUser.UserPasswordReset,
		dbUser.UserCreatedAt,
	)

	return user, nil
}

func (userDatabase *UserDatabase) FindUserById(id string) (*userDomain.User, error) {
	var dbUser *models.UserWithPermissionGroup

	query := `SELECT
        id as user_id,
        name as user_name,
        email as user_email,
		password as user_password,
		cpf as user_cpf,
        role as user_role,
		cellphone_number as user_cellphone_number,
        is_active as user_is_active,
        status as user_status,
        created_at as user_created_at,
      FROM user
      WHERE id = ?`

	if err := userDatabase.connection.Raw(query, &dbUser, id); err != nil {
		return nil, err
	}

	if dbUser == nil {
		return nil, nil
	}

	user := userDomain.NewUser(
		dbUser.UserID,
		dbUser.UserName,
		dbUser.UserEmail,
		dbUser.UserPassword,
		dbUser.UserCpf,
		dbUser.UserCellphoneNumber,
		dbUser.UserStatus,
		dbUser.UserIsActive,
		dbUser.UserPasswordReset,
		dbUser.UserCreatedAt,
	)

	return user, nil
}

func (userDatabase *UserDatabase) CreateUser(user *userDomain.User) error {
	createdAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}
	updatedAt := sql.NullTime{Valid: false}
	deletedAt := sql.NullTime{Valid: false}

	query := `INSERT INTO user (
			   		id,
					name,
			 		email,
		 			password,
					cpf,
					cellphone_number,
					status,
			 		created_at,
			 		updated_at,
			 		deleted_at
				) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var statement interface{}

	if err := userDatabase.connection.Raw(
		query,
		statement,
		user.Id,
		strings.ToLower(user.Name),
		strings.ToLower(user.Email),
		user.Password,
		user.Cpf,
		user.CellphoneNumber,
		user.Status,
		createdAt,
		updatedAt,
		deletedAt,
	); err != nil {
		return err
	}

	return nil
}

func (userDatabase *UserDatabase) UpdateResetPasswordStatus(toggle bool, status string, userId string) error {
	query := `UPDATE user SET password_reset = ?, status = ? WHERE id = ?`

	var userDb models.User

	if err := userDatabase.connection.Raw(query, &userDb, toggle, status, userId); err != nil {
		return err
	}

	return nil
}
