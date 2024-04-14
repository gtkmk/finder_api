package like

import (
	"os"

	"github.com/gtkmk/finder_api/adapter/http/handlers/likeHandler"
	"github.com/gtkmk/finder_api/adapter/http/middleware"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/jwtAuth"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/envMode"
	"github.com/gtkmk/finder_api/infra/httpContextValuesExtractor"

	"github.com/gin-gonic/gin"
)

const (
	CreateLikeCreateLikeConst string = "CreateLike"
	// === Route constants marker ===
)

type LikeRoutes struct {
	*gin.Engine
	likeHandlers map[string]port.HandlerInterface
	jwt          *middleware.IsAuthorized
}

func NewLikeRoutes(
	app *gin.Engine,
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	passwordEncryption port.EncryptionInterface,
	notificationService port.NotificationInterface,
) port.RoutesInterface {
	jwt := jwtAuth.NewjwtAuth(os.Getenv(envMode.JwtSecretConst))

	return &LikeRoutes{
		app,
		createMapOfLikeHandlers(connection, notificationService, uuid, passwordEncryption),
		middleware.NewIsAuthorized(
			jwt,
			connection,
			uuid,
		),
	}
}

func (likeRoutes *LikeRoutes) Register() {
	likeRoutes.POST(
		routesConstants.PostLikeRouteConst,
		likeRoutes.jwt.IsAuthorizedMiddleware(),
		likeRoutes.likeHandlers[CreateLikeCreateLikeConst].Handle,
	)
	// === Register route marker ===
}

func createMapOfLikeHandlers(
	connection port.ConnectionInterface,
	notificationService port.NotificationInterface,
	uuid port.UuidInterface,
	passwordEncryption port.EncryptionInterface,
) map[string]port.HandlerInterface {
	contextExtractor := httpContextValuesExtractor.NewHttpContextValuesExtractor()

	return map[string]port.HandlerInterface{
		CreateLikeCreateLikeConst: likeHandler.NewCreateLikeCreateLikeHandler(
			connection,
			uuid,
			contextExtractor,
		),
		// === Register handler marker ===
	}
}
