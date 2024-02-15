package sharedMethods

type CreateUserEventInterface interface {
	SaveLoginEvent(userId string, ip string, device string) error
	SaveNewUserEvent(userId string, ip string, device string, newValue string) error
}
