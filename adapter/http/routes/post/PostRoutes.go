package user

import (
	"github.com/gtkmk/finder_api/adapter/http/handlers/postHandler"
	"github.com/gtkmk/finder_api/adapter/http/middleware"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/jwtAuth"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/envMode"
	"github.com/gtkmk/finder_api/infra/httpContextValuesExtractor"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	PostPostKeyConst string = "postPost"
)

type PostRoutes struct {
	*gin.Engine
	postHandlers map[string]port.HandlerInterface
	jwt          *middleware.IsAuthorized
}

func NewPostRoutes(
	app *gin.Engine,
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	passwordEncryption port.EncryptionInterface,
	notificationService port.NotificationInterface,
) port.RoutesInterface {
	jwt := jwtAuth.NewjwtAuth(os.Getenv(envMode.JwtSecretConst))

	return &PostRoutes{
		app,
		createMapOfUserHandlers(connection, notificationService, uuid, passwordEncryption),
		middleware.NewIsAuthorized(
			jwt,
			connection,
			uuid,
		),
	}
}

func (postRoutes *PostRoutes) Register() {
	postRoutes.POST(
		routesConstants.PostCreatePostRouteCOnst,
		postRoutes.jwt.IsAuthorizedMiddleware(),
		postRoutes.postHandlers[PostPostKeyConst].Handle,
	)
}

func createMapOfUserHandlers(
	connection port.ConnectionInterface,
	notificationService port.NotificationInterface,
	uuid port.UuidInterface,
	passwordEncryption port.EncryptionInterface,
) map[string]port.HandlerInterface {
	contextExtractor := httpContextValuesExtractor.NewHttpContextValuesExtractor()

	return map[string]port.HandlerInterface{
		PostPostKeyConst: postHandler.NewCreatePostHandler(
			connection,
			uuid,
		),
	}
}
