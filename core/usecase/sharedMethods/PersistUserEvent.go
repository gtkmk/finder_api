package sharedMethods

import (
	"github.com/gtkmk/finder_api/core/domain/userEventDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	"github.com/gtkmk/finder_api/core/port/sharedMethods"
)

type PersistUserEvent struct {
	userEventDatabase repositories.UserEventRepositoryInterface
	uuid              port.UuidInterface
	customError       port.CustomErrorInterface
}

func NewPersistUserEvent(
	userEventDatabase repositories.UserEventRepositoryInterface,
	uuid port.UuidInterface,
	customErrorInterface port.CustomErrorInterface,
) sharedMethods.CreateUserEventInterface {
	return &PersistUserEvent{
		userEventDatabase: userEventDatabase,
		uuid:              uuid,
		customError:       customErrorInterface,
	}
}

func (persistUserEvent *PersistUserEvent) SaveLoginEvent(userId string, ip string, device string) error {
	userEvent := userEventDomain.NewUserEvent(
		persistUserEvent.uuid.GenerateUuid(),
		userId,
		userEventDomain.LoginEventConst,
		nil,
		nil,
		ip,
		device,
	)

	if err := persistUserEvent.userEventDatabase.Save(userEvent); err != nil {
		return persistUserEvent.customError.ThrowError(err.Error())
	}

	return nil
}

func (persistUserEvent *PersistUserEvent) SaveNewUserEvent(userId string, ip string, device string, newValue string) error {
	userEvent := userEventDomain.NewUserEvent(
		persistUserEvent.uuid.GenerateUuid(),
		userId,
		userEventDomain.NewUserEventConst,
		nil,
		&newValue,
		ip,
		device,
	)

	if err := persistUserEvent.userEventDatabase.Save(userEvent); err != nil {
		return persistUserEvent.customError.ThrowError(err.Error())
	}

	return nil
}
