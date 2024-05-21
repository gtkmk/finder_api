package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/envMode"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBFactory struct {
	DbUserName string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

func NewDBFactory() port.DbFactoryInterface {
	dbUserName := os.Getenv(envMode.DataBaseUsernameConst)
	dbPassword := os.Getenv(envMode.DataBasePasswordConst)
	dbHost := os.Getenv(envMode.DataBaseHostConst)
	dbPort := os.Getenv(envMode.DataBasePortConst)
	dbName := os.Getenv(envMode.DataBaseNameConst)
	return &DBFactory{
		DbUserName: dbUserName,
		DbPassword: dbPassword,
		DbHost:     dbHost,
		DbPort:     dbPort,
		DbName:     dbName,
	}
}

func (dbFactory *DBFactory) CreateDBConnection() (*gorm.DB, error) {
	isTestEnv := os.Getenv(envMode.TestEnvConst) == "true"

	if isTestEnv {
		return dbFactory.createSqliteConnection()
	}

	return dbFactory.createSqlConnection()
}

func (dbFactory *DBFactory) createSqliteConnection() (*gorm.DB, error) {
	rootDir, err := dbFactory.getRootDir()
	if err != nil {
		panic("Error getting root directory")
	}

	dbPath := filepath.Join(rootDir, "gorm.db")
	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}

func (dbFactory *DBFactory) getRootDir() (string, error) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	return basePath, nil
}

func (dbFactory *DBFactory) createSqlConnection() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: false,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,        // Don't include params in the SQL log
			Colorful:                  true,         // Disable color
		},
	)

	gormConfig := &gorm.Config{Logger: newLogger, SkipDefaultTransaction: true}

	dnsString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbFactory.DbUserName,
		dbFactory.DbPassword,
		dbFactory.DbHost,
		dbFactory.DbPort,
		dbFactory.DbName,
	)

	sqlConfig := mysql.Config{DSN: dnsString}

	db, err := gorm.Open(mysql.New(sqlConfig), gormConfig)

	if err != nil {
		return nil, err
	}

	return db, err
}
