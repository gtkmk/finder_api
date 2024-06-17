package userRequestEntity

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UpdateUserInfoRequest struct {
	Name            string `json:"name"`
	CellphoneNumber string `json:"cellphone_number"`
}

func NewUpdateUserInfoRequest(
	context *gin.Context,
) (*UpdateUserInfoRequest, error) {
	updateUserInfoRequest := &UpdateUserInfoRequest{}
	if err := context.ShouldBind(updateUserInfoRequest); err != nil {
		return nil, err
	}

	return updateUserInfoRequest, nil
}

func (updateUserInfoRequest *UpdateUserInfoRequest) Validate() error {
	nameValidationError := ValidateName(updateUserInfoRequest.Name)

	if nameValidationError != nil {
		return nameValidationError
	}

	if updateUserInfoRequest.CellphoneNumber != "" {
		cellphoneNumberValidationError := ValidateCellphoneNumber(updateUserInfoRequest.CellphoneNumber)

		if cellphoneNumberValidationError != nil {
			return cellphoneNumberValidationError
		}
	}

	return nil
}

func (updateUserInfoRequest *UpdateUserInfoRequest) DecodeUpdatedUserInfoRequest(req *http.Request) (*UpdateUserInfoRequest, error) {
	updateUserInfoRequest.Name = strings.ToLower(updateUserInfoRequest.Name)

	return updateUserInfoRequest, nil
}
