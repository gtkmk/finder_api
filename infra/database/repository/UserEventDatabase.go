package repository

import (
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/userEventDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type UserEventRepository struct {
	connection port.ConnectionInterface
}

func NewUserEventRepository(connection port.ConnectionInterface) repositories.UserEventRepositoryInterface {
	return &UserEventRepository{
		connection: connection,
	}
}

func (userEventRepository *UserEventRepository) Save(userEvent *userEventDomain.UserEvent) error {
	createdAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}

	query := `INSERT INTO user_event (
				  id,
				  user_id,
				  event,
				  old_value,
				  new_value,
				  ip,
				  device,
				  created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	var statement interface{}

	if err := userEventRepository.connection.Raw(
		query,
		&statement,
		userEvent.Id,
		userEvent.UserId,
		userEvent.Event,
		userEvent.OldValue,
		userEvent.NewValue,
		userEvent.Ip,
		userEvent.Device,
		createdAt,
	); err != nil {
		return err
	}

	return nil
}
