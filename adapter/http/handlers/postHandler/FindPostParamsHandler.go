package postHandler

import (
	"github.com/gin-gonic/gin"

	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/port"
	postUsecase "github.com/gtkmk/finder_api/core/usecase/post"
)

type FindPostPostParamsHandler struct {
	connection port.ConnectionInterface
	uuid       port.UuidInterface
}

func NewFindPostPostParamsHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
) port.HandlerInterface {
	return &FindPostPostParamsHandler{
		connection: connection,
		uuid:       uuid,
	}
}

func (findPostPostParamsHandler *FindPostPostParamsHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, findPostPostParamsHandler.connection, findPostPostParamsHandler.uuid)

	postParams := postUsecase.NewFindPostPostParams().Execute()

	jsonResponse.SendJson("postParams", postParams, routesConstants.StatusOk)
}
