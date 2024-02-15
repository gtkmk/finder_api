package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID              string `gorm:"primaryKey;type:varchar(191);"`
	Name            string `gorm:"not null;type:varchar(100);" json:"name"`
	Email           string `gorm:"unique;not null;type:varchar(100);" json:"email"`
	Password        string `gorm:"not null;type:varchar(100);" json:"password"`
	Cpf             string `gorm:"unique;not null;type:varchar(14);" json:"cpf"`
	CellphoneNumber string `gorm:"type:varchar(15);" json:"cellphone_number"`
	IsActive        bool   `gorm:"default:true;not null;type:tinyint(1);" json:"is_active"`
	Status          string `gorm:"column:status;not null;type:enum('expired', 'pending', 'logged');"`
	PasswordReset   bool   `gorm:"default:false;"`
	CreatedAt       time.Time
	UpdatedAt       sql.NullTime
	DeletedAt       sql.NullTime
}

type UserWithPermissionGroup struct {
	UserID              string
	UserName            string
	UserEmail           string
	UserPassword        string
	UserCpf             string
	UserCellphoneNumber string
	UserStatus          string
	UserIsActive        bool
	UserCreatedAt       time.Time
	UserPasswordReset   bool
}
