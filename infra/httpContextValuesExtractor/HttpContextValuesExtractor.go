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
	loggedUserLayer float64,
	extractError error,
) {
	userId, exists := context.Get("userId")

	if !exists {
		return "", 0, helper.ErrorBuilder(helper.ErrorExtractingUserIdConst)
	}

	groupLayer, exists := context.Get("groupLayer")

	if !exists {
		return "", 0, helper.ErrorBuilder(helper.ErrorExtractingUserGroupConst)
	}

	return userId.(string), groupLayer.(float64), nil
}
