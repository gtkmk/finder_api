package notification

import (
	"github.com/gtkmk/finder_api/core/domain/notificationDomain"
	"github.com/gtkmk/finder_api/core/port"
)

type Notification struct {
	httpClient port.ClientInterface
}

func NewNotification(httpClient port.ClientInterface) port.NotificationInterface {
	return &Notification{httpClient}
}

func (notify *Notification) SendNotifications(notification *notificationDomain.Notification) error {
	// TODO: Implementar notificação
	return nil
}
