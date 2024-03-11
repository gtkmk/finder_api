package httpContextValuesExtractor

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/core/domain/helper"

	"github.com/gtkmk/finder_api/core/port"
)

type HttpContextValuesExtractor struct{}

func NewHttpContextValuesExtractor() port.HttpContextValuesExtractorInterface {
	return &HttpContextValuesExtractor{}
}

func (httpContextValuesExtractor *HttpContextValuesExtractor) Extract(context *gin.Context) (
	loggedUserId string,
	extractError error,
) {
	userId, exists := context.Get("userId")

	if !exists {
		return "", helper.ErrorBuilder(helper.ErrorExtractingUserIdConst)
	}

	return userId.(string), nil
}
