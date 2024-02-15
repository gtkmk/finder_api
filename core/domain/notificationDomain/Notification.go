package notificationDomain

import (
	emailDomain "github.com/gtkmk/finder_api/core/domain/email"
	"github.com/gtkmk/finder_api/core/domain/postDomain"
)

type Notification struct {
	NotificationType string
	Name             string
	Email            string
	Cellphone        string
	Cpf              string
	Url              string
	Post             *postDomain.Post
}

func NewNotification(
	notificationType string,
	name string,
	email string,
	cellphone string,
	cpf string,
	url string,
	post *postDomain.Post,
) *Notification {
	email = emailDomain.NewEmail(email).AsString()

	return &Notification{
		notificationType,
		name,
		email,
		cellphone,
		cpf,
		url,
		post,
	}
}
