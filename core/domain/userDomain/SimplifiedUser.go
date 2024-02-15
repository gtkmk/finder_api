package userDomain

type SimplifiedUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewSimplifiedUser(
	id string,
	name string,
) *SimplifiedUser {
	return &SimplifiedUser{
		Id:   id,
		Name: name,
	}
}
