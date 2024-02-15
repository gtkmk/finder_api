package port

import "github.com/gtkmk/finder_api/core/domain/notificationDomain"

type NotificationInterface interface {
	SendNotifications(notification *notificationDomain.Notification) error
}
