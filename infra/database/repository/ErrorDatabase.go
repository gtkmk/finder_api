package repository

import (
	"log"

	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type ErrorDatabase struct {
	connection port.ConnectionInterface
}

func NewErrorDatabase(connection port.ConnectionInterface) repositories.ErrorRepositoryInterface {
	return &ErrorDatabase{
		connection,
	}
}

func (errorDatabase ErrorDatabase) SaveErrorWithoutStack(id string, errorMessage string, createdAt string) {
	query := `INSERT INTO error (
			id,
        	message,
            created_at
        ) VALUES (?, ?, ?)`

	var statement interface{}

	if err := errorDatabase.connection.Raw(
		query,
		statement,
		id,
		errorMessage,
		createdAt,
	); err != nil {
		log.Println("Failed to save error:", errorMessage)
	}
}

func (errorDatabase ErrorDatabase) SaveErrorWithStack(id string, errorMessage string, stack string, createdAt string) {
	query := `INSERT INTO error (
			id,
			message,
			stack,
            created_at
        ) VALUES (?, ?, ?, ?)`

	var statement interface{}

	if err := errorDatabase.connection.Raw(
		query,
		statement,
		id,
		errorMessage,
		stack,
		createdAt,
	); err != nil {
		log.Println("Failed to save error:", errorMessage)
	}
}
