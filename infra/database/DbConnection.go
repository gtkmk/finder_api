package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/envMode"
	"gorm.io/gorm"
)

type DBConnection struct {
	connection *gorm.DB
}

func NewDBConnection() port.ConnectionInterface {
	return &DBConnection{}
}

const (
	TransactionCommitErrorConst     = "TransactionCommitError"
	TransactionRollbackErrorConst   = "TransactionRollbackError"
	TransactionSavePointErrorConst  = "TransactionSavePointError"
	TransactionRollbackToErrorConst = "TransactionRollbackToError"
	ExecuteRowsQueryErrorConst      = "ExecuteRowsQueryError"
	ExecuteRawStatementErrorConst   = "ExecuteRawStatementError"
)

func (db *DBConnection) Open() error {
	var database *gorm.DB
	var err error

	maxRetries := 10
	retryInterval := 5 * time.Second

	dbFactory := NewDBFactory()

	for retryCount := 1; retryCount <= maxRetries; retryCount++ {
		database, err = dbFactory.CreateDBConnection()
		if err == nil {
			break
		}

		log.Printf("Failed to connect to the database. Retrying in %v...", retryInterval)
		time.Sleep(retryInterval)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	sqlDB, ok := database.ConnPool.(*sql.DB)
	if !ok {
		return fmt.Errorf("unexpected connection pool type")
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	db.connection = database
	log.Printf("Database connection successfully opened")
	return nil
}

func (db *DBConnection) Raw(query string, statment interface{}, values ...any) error {
	transaction := db.connection.Raw(query, values...).Scan(statment)

	if transaction.Error != nil {
		return db.errorBuilder(ExecuteRawStatementErrorConst, transaction.Error)
	}

	return nil
}

func (db *DBConnection) Rows(query string, values ...any) ([]map[string]interface{}, error) {
	rows, err := db.connection.Raw(query, values...).Rows()
	if err != nil {
		return nil, db.errorBuilder(ExecuteRowsQueryErrorConst, err)
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	count := len(columns)
	statementPtrs := make([]interface{}, count)

	for i := range statementPtrs {
		var v interface{}
		statementPtrs[i] = &v
	}

	result := make([]map[string]interface{}, 0)

	for rows.Next() {
		err := rows.Scan(statementPtrs...)
		if err != nil {
			continue
		}
		row := make(map[string]interface{})
		for i, col := range columns {
			val := *(statementPtrs[i].(*interface{}))
			row[col] = handleScanReturn(val)
		}
		result = append(result, row)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return result, nil
}

func handleScanReturn(result any) any {
	switch value := result.(type) {
	case bool:
		return value
	case int64:
		return value
	case float32:
		return value
	case float64:
		return value
	case int:
		return value
	case string:
		return value
	case []byte:
		return string(value)
	case time.Time:
		return value.Format(datetimeDomain.DATETIME_FORMAT)
	case nil:
		return nil
	default:
		fmt.Println(helper.UnknownResultTypeConst)
		return nil
	}
}

func (db *DBConnection) Close() error {
	sqlDB, err := db.connection.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("Falha ao fechar conexao com o banco: %w", err)
	}

	log.Println("Database connection successfully closed")

	return nil
}

func (db *DBConnection) BeginTransaction() (port.ConnectionInterface, error) {
	tx := db.connection.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &DBConnection{tx}, nil
}

func (db *DBConnection) Commit() error {
	transaction := db.connection.Commit()

	if transaction.Error != nil {
		return db.errorBuilder(TransactionCommitErrorConst, transaction.Error)
	}

	return nil
}

func (db *DBConnection) Rollback() error {
	transaction := db.connection.Rollback()

	if transaction.Error != nil {
		return db.errorBuilder(TransactionRollbackErrorConst, transaction.Error)
	}

	return nil
}

func (db *DBConnection) SavePoint(checkPointName string) error {
	transaction := db.connection.SavePoint(checkPointName)

	if transaction.Error != nil {
		return db.errorBuilder(TransactionSavePointErrorConst, transaction.Error)
	}

	return nil
}

func (db *DBConnection) RollbackTo(checkPointName string) error {
	transaction := db.connection.RollbackTo(checkPointName)

	if transaction.Error != nil {
		return db.errorBuilder(TransactionRollbackToErrorConst, transaction.Error)
	}

	return nil
}

func (db *DBConnection) errorBuilder(errorType string, err error) error {
	isProd := os.Getenv(envMode.IsProdConst)

	if isProd == "1" {
		return db.errorStringBuilder(errorType)
	}

	return err
}

func (db *DBConnection) errorStringBuilder(errorType string) error {
	var err error
	switch errorType {
	case TransactionCommitErrorConst:
		err = helper.ErrorBuilder(helper.ErrorWithCodeConst, helper.ErrorWhenTryToCommitTransactionCodeConst)
	case TransactionRollbackErrorConst:
		err = helper.ErrorBuilder(helper.ErrorWithCodeConst, helper.ErrorWhenTryToRollbackTransactionCodeConst)
	case TransactionSavePointErrorConst:
		err = helper.ErrorBuilder(helper.ErrorWithCodeConst, helper.ErrorWhenTryToCreateATransactionSavePointCodeConst)
	case TransactionRollbackToErrorConst:
		err = helper.ErrorBuilder(helper.ErrorWithCodeConst, helper.ErrorWhenTryToRollbackToSavePointCodeConst)
	case ExecuteRowsQueryErrorConst:
		err = helper.ErrorBuilder(helper.ErrorWithCodeConst, helper.ErrorWhenTryToExecuteRowsQueryCodeConst)
	case ExecuteRawStatementErrorConst:
		err = helper.ErrorBuilder(helper.ErrorWithCodeConst, helper.ErrorWhenExecuteRawStatementCodeConst)
	}

	return err
}
