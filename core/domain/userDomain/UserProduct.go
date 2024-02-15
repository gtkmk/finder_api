package userDomain

type UserProducts struct {
	ID        string `json:"id"`
	UserId    string `json:"user_id"`
	ProductId string `json:"product_id"`
}

func NewUserProducts(id string, userId string, productId string) *UserProducts {
	return &UserProducts{
		id,
		userId,
		productId,
	}
}
