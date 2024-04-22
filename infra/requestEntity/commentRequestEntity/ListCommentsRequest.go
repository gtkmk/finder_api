package commentRequestEntity

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

type ListCommentsRequest struct {
	uuid   port.UuidInterface
	PostID string `form:"post_id"`
	Page   *int64 `form:"page"`
}

const (
	PostIdFieldConst = "a postagem"
)

func NewListCommentsRequest(
	context *gin.Context,
	uuid port.UuidInterface,
) (*ListCommentsRequest, error) {
	listCommentsRequest := &ListCommentsRequest{}
	if err := context.ShouldBind(listCommentsRequest); err != nil {
		return nil, err
	}

	return listCommentsRequest, nil
}

func (listCommentsRequest *ListCommentsRequest) ValidateCommentsFilterFields(context *gin.Context) error {
	if err := listCommentsRequest.verifyIfPageIsValid(listCommentsRequest.Page); err != nil {
		return err
	}

	if err := requestEntityFieldsValidation.IsValidUUID(PostIdFieldConst, listCommentsRequest.PostID); err != nil {
		return err
	}

	return nil
}

func (listCommentsRequest *ListCommentsRequest) verifyIfPageIsValid(page *int64) error {
	if page == nil || *page <= 0 {
		return fmt.Errorf(
			helper.FieldIsMandatoryAndMustToBeGreaterThanZeroConst,
			"página",
		)
	}
	return nil
}

func (listCommentsRequest *ListCommentsRequest) RetrieveCommentsFiltersInfo() (string, int64) {
	return listCommentsRequest.PostID, *listCommentsRequest.Page
}
