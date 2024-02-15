package userEventDomain

const (
	LoginEventConst    = "login"
	NewUserEventConst  = "new_user"
	EditUserEventConst = "edit_user"
)

type UserEvent struct {
	Id       string
	UserId   string
	Event    string
	OldValue *string
	NewValue *string
	Ip       string
	Device   string
}

func NewUserEvent(
	id string,
	userId string,
	event string,
	oldValue *string,
	newValue *string,
	ip string,
	device string,
) *UserEvent {
	return &UserEvent{
		Id:       id,
		UserId:   userId,
		Event:    event,
		OldValue: oldValue,
		NewValue: newValue,
		Ip:       ip,
		Device:   device,
	}
}
