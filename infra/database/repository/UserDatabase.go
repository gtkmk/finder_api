package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/userDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	"github.com/gtkmk/finder_api/infra/database/models"
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

func (userDatabase *UserDatabase) VerifyIfUserExistsByUserName(userName string) bool {
	var exists bool
	query := `SELECT CASE WHEN EXISTS (SELECT 1 FROM user WHERE user_name = ?) THEN 1 ELSE 0 END`

	if err := userDatabase.connection.Raw(query, &exists, userName); err != nil {
		return false
	}

	return exists
}

func (userDatabase *UserDatabase) FindUserByEmail(email string) (*userDomain.User, error) {
	var dbUser *models.UserWithPermissionGroup

	query := `SELECT
        id as user_id,
        name,
		user_name,
        email as user_email,
		password as user_password,
		cpf as user_cpf,
		cellphone_number as user_cellphone_number,
		status as user_status,
        is_active as user_is_active,
        created_at as user_created_at
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
		dbUser.Name,
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
        name,
		user_name,
        email as user_email,
		password as user_password,
		cpf as user_cpf,
		cellphone_number as user_cellphone_number,
        is_active as user_is_active,
        status as user_status,
        created_at as user_created_at
      FROM 
	  	user
      WHERE
	  	id = ? AND deleted_at IS NULL
	`

	if err := userDatabase.connection.Raw(query, &dbUser, id); err != nil {
		return nil, err
	}

	if dbUser == nil {
		return nil, nil
	}

	user := userDomain.NewUser(
		dbUser.UserID,
		dbUser.Name,
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
					user_name,
			 		email,
		 			password,
					cpf,
					cellphone_number,
					status,
			 		created_at,
			 		updated_at,
			 		deleted_at
				) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var statement interface{}

	if err := userDatabase.connection.Raw(
		query,
		statement,
		user.Id,
		strings.ToLower(user.Name),
		user.UserName,
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

func (userDatabase *UserDatabase) SetUserStatus(userId string, status string) error {
	query := `UPDATE user SET status = ? WHERE id = ?`
	var userDb models.User

	if err := userDatabase.connection.Raw(query, &userDb, status, userId); err != nil {
		return err
	}

	return nil
}

func (userDatabase *UserDatabase) ResetUserPassword(
	userId string,
	password string,
) error {
	updatedAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}

	query := `UPDATE user
                    SET password = ?, updated_at = ?  
                    WHERE id = ?`

	var userDb models.User

	if err := userDatabase.connection.Raw(query, &userDb, password, updatedAt, userId); err != nil {
		return err
	}

	return nil
}

func (userDatabase *UserDatabase) FindCompleteUserInfoByID(userId, loggedUserId string) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			user.id,
			user.name,
			user.user_name,
			user.email,
			user.cellphone_number,
			user.status,
			user.is_active,
			user.created_at,
			(SELECT COUNT(*) FROM follow WHERE follower_id = user.id) AS following_count,
			(SELECT COUNT(*) FROM follow WHERE followed_id = user.id) AS followers_count,
			(SELECT COUNT(*) FROM post WHERE user_id = user.id AND lost_found = 'lost' AND deleted_at IS NULL) AS lost_posts_count,
			(SELECT COUNT(*) FROM post WHERE user_id = user.id AND lost_found = 'found' AND deleted_at IS NULL) AS found_posts_count,
			(SELECT COUNT(*) FROM post WHERE user_id = user.id AND deleted_at IS NULL) AS total_posts_count,
			docPic.path AS profile_picture_path,
			docBan.path AS profile_banner_picture_path,
			CASE
				WHEN user.id = ? THEN true
				ELSE false
			END AS is_own_profile,
			CASE
				WHEN EXISTS (SELECT 1 FROM follow WHERE follower_id = ? AND followed_id = user.id) THEN true
				ELSE false
			END AS is_following,
			CASE
				WHEN EXISTS (SELECT 1 FROM follow WHERE follower_id = user.id AND followed_id = ?) THEN true
				ELSE false
			END AS is_followed
		FROM
				user
		INNER JOIN 
				document docPic ON user.id = docPic.owner_id AND docPic.type = 'profile_picture' AND docPic.deleted_at IS NULL
		INNER JOIN
				document docBan ON user.id = docBan.owner_id AND docBan.type = 'profile_banner_picture' AND docBan.deleted_at IS NULL
		WHERE 
    		user.id = ? AND user.deleted_at IS NULL
	`

	dbUser, err := userDatabase.connection.Rows(
		query,
		loggedUserId,
		loggedUserId,
		loggedUserId,
		userId,
	)

	if err != nil {
		return nil, err
	}

	return dbUser, nil
}

func (userDatabase *UserDatabase) FindUsersListByName(userName string, loggedUserId string) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			user.id,
			user.name,
			user.user_name,
			CASE
				WHEN EXISTS (SELECT 1 FROM follow WHERE follower_id = ? AND followed_id = user.id) THEN true
				ELSE false
			END AS is_following,
			CASE
				WHEN EXISTS (SELECT 1 FROM follow WHERE follower_id = user.id AND followed_id = ?) THEN true
				ELSE false
			END AS is_followed,
			document.path AS profilePicture
		FROM
			user
			LEFT JOIN
			document ON document.owner_id = user.id AND document.type = 'profile_picture' AND document.deleted_at IS NULL
		WHERE 
			(user.name LIKE ? OR user.user_name LIKE ?)
			AND user.id != ? 
			AND user.deleted_at IS NULL
	`

	return userDatabase.connection.Rows(
		query,
		loggedUserId,
		loggedUserId,
		fmt.Sprintf("%%%s%%", userName),
		fmt.Sprintf("%%%s%%", userName),
		loggedUserId,
	)
}
