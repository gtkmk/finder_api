package models

import (
	"database/sql"
	"time"
)

type Document struct {
	ID            string `gorm:"primaryKey;type:varchar(191);"`
	Type          string `gorm:"not null;type:enum('media', 'profile_picture', 'profile_banner_picture');"`
	Path          string `gorm:"not null;type:varchar(250);"`
	MimeType      string `gorm:"not null;type:varchar(100);"`
	PostId        string `gorm:"type:varchar(191);"`
	OwnerId       string `gorm:"not null;type:varchar(191);"`
	DeletedReason string `gorm:"type:varchar(100);"`
	CreatedAt     time.Time
	UpdatedAt     sql.NullTime
	DeletedAt     sql.NullTime
}
