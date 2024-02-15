package middleware

import (
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/jwtAuth"
	"github.com/gtkmk/finder_api/core/port"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IsAuthorized struct {
	jwt        *jwtAuth.JwtAuth
	handler    gin.HandlerFunc
	connection port.ConnectionInterface
	uuid       port.UuidInterface
}

func NewIsAuthorized(
	jwt *jwtAuth.JwtAuth,
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
) *IsAuthorized {
	return &IsAuthorized{
		jwt,
		nil,
		connection,
		uuid,
	}
}

func (isAuthorized *IsAuthorized) IsAuthorizedMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		jsonResponse := routes.NewJsonResponse(context, isAuthorized.connection, isAuthorized.uuid)

		_, err := context.Request.Cookie("token")

		if err != nil {
			jsonResponse.SendJson(
				routesConstants.MessageKeyConst,
				err.Error(),
				routesConstants.Unauthorized,
			)

			context.Abort()

			return
		}

		token, err := isAuthorized.jwt.CheckJwt(context.Request)

		if err != nil {
			jsonResponse.SendJson(
				routesConstants.MessageKeyConst,
				"Unauthorized",
				routesConstants.Unauthorized,
			)

			context.Abort()

			return
		}

		if token == nil {
			jsonResponse.SendJson(
				routesConstants.MessageKeyConst,
				"Unauthorized",
				routesConstants.Unauthorized,
			)

			context.JSON(http.StatusUnauthorized, err)
			context.Abort()

			return
		}

		context.Set("userId", (*token)["i"])
		context.Set("groupLayer", (*token)["l"])
		context.Next()
	}
}
