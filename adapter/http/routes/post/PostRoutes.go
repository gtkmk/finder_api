package post

import (
	"os"

	"github.com/gtkmk/finder_api/adapter/http/handlers/postHandler"
	"github.com/gtkmk/finder_api/adapter/http/middleware"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/jwtAuth"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/envMode"
	"github.com/gtkmk/finder_api/infra/httpContextValuesExtractor"

	"github.com/gin-gonic/gin"
)

const (
	CreatePost   string = "createPost"
	FindAllPosts string = "findAllPosts"
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
		routesConstants.PostCreatePostRouteConst,
		postRoutes.jwt.IsAuthorizedMiddleware(),
		postRoutes.postHandlers[CreatePost].Handle,
	)

	postRoutes.GET(
		routesConstants.PostFindAllPostsRouteConst,
		postRoutes.jwt.IsAuthorizedMiddleware(),
		postRoutes.postHandlers[FindAllPosts].Handle,
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
		CreatePost: postHandler.NewCreatePostHandler(
			connection,
			uuid,
			contextExtractor,
		),
		FindAllPosts: postHandler.NewFindPostAllHandler(
			connection,
			uuid,
			contextExtractor,
		),
	}
}
