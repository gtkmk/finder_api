package port

import "gorm.io/gorm"

type DbFactoryInterface interface {
	CreateDBConnection() (*gorm.DB, error)
}
