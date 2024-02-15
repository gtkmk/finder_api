package signHandler

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/envMode"
)

type SignOutHandler struct {
	connection port.ConnectionInterface
	uuid       port.UuidInterface
}

func NewSignOutHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
) port.HandlerInterface {
	return &SignOutHandler{
		connection,
		uuid,
	}
}

func (signOutHandler *SignOutHandler) Handle(context *gin.Context) {
	context.SetCookie(
		"token",
		"",
		-1, // invalidating cookie
		"/",
		os.Getenv(envMode.ApplicationDomainConst),
		true,
		true,
	)

	context.JSON(routesConstants.StatusOk, map[string]string{routesConstants.MessageKeyConst: "Deslogado com sucesso"})
}
