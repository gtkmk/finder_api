package userDomain

type ManagerOrConsultant struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewManager(
	Id string,
	Name string,
) *ManagerOrConsultant {
	return &ManagerOrConsultant{
		Id,
		Name,
	}
}
