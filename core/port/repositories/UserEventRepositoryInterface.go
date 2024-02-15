package repositories

import "github.com/gtkmk/finder_api/core/domain/userEventDomain"

type UserEventRepositoryInterface interface {
	Save(userEvent *userEventDomain.UserEvent) error
}
