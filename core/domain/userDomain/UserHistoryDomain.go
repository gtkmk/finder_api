package userDomain

type UserHistory struct {
	Id          string `json:"id"`
	Responsible string `json:"responsible"`
	Email       string `json:"email"`
	Cellphone   string `json:"cellphone"`
	HistoryId   string `json:"history_id"`
}

func NewUserHistory(
	id string,
	responsible string,
	email string,
	cellphone string,
	historyId string,
) *UserHistory {
	return &UserHistory{
		id,
		responsible,
		email,
		cellphone,
		historyId,
	}
}
