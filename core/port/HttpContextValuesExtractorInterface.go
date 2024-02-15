package port

import (
	"github.com/gin-gonic/gin"
)

type Permission struct {
	I  string `json:"i"`
	RN string `json:"rn"`
	OP string `json:"op"`
}

type PermissionJWT struct {
	Expiration int64         `json:"expiration"`
	P          []*Permission `json:"p"`
}

type HttpContextValuesExtractorInterface interface {
	Extract(context *gin.Context) (
		loggedUserId string,
		loggedUserLayer float64,
		extractError error,
	)
}
