package userRequestEntity

import (
	"encoding/json"
	"net/http"

	"github.com/gtkmk/finder_api/infra/EmptyBodyVerification"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

type UserProductsRequest struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	ProductId string `json:"product_id"`
}

const (
	ProductIdFieldConst = "o produto"
)

func DecodeUserProductRequest(req *http.Request) (*UserProductsRequest, error) {
	if err := EmptyBodyVerification.ValidateBody(req); err != nil {
		return nil, err
	}

	var userProduct *UserProductsRequest
	if err := json.NewDecoder(req.Body).Decode(&userProduct); err != nil {
		return nil, err
	}

	return userProduct, nil
}

func (user *UserProductsRequest) Validate() error {
	if err := requestEntityFieldsValidation.IsValidUUID(UserIdFieldConst, user.UserId); err != nil {
		return err
	}

	return requestEntityFieldsValidation.IsValidUUID(ProductIdFieldConst, user.ProductId)
}
