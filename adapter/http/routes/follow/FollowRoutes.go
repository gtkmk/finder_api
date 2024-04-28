package follow

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/handlers/followHandler"
	"github.com/gtkmk/finder_api/adapter/http/middleware"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/jwtAuth"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/envMode"
	"github.com/gtkmk/finder_api/infra/httpContextValuesExtractor"
)

const (
	CreateFollowUserConst string = "CreateFollowUser"
	// === Route constants marker ===
)

type FollowRoutes struct {
	*gin.Engine
	likeHandlers map[string]port.HandlerInterface
	jwt          *middleware.IsAuthorized
}

func NewFollowRoutes(
	app *gin.Engine,
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	passwordEncryption port.EncryptionInterface,
	notificationService port.NotificationInterface,
) port.RoutesInterface {
	jwt := jwtAuth.NewjwtAuth(os.Getenv(envMode.JwtSecretConst))

	return &FollowRoutes{
		app,
		createMapOfFollowHandlers(connection, notificationService, uuid, passwordEncryption),
		middleware.NewIsAuthorized(
			jwt,
			connection,
			uuid,
		),
	}
}

func (followRoutes *FollowRoutes) Register() {
	followRoutes.POST(
		routesConstants.PostFollowRouteConst,
		followRoutes.jwt.IsAuthorizedMiddleware(),
		followRoutes.likeHandlers[CreateFollowUserConst].Handle,
	)
	// === Register route marker ===
}

func createMapOfFollowHandlers(
	connection port.ConnectionInterface,
	notificationService port.NotificationInterface,
	uuid port.UuidInterface,
	passwordEncryption port.EncryptionInterface,
) map[string]port.HandlerInterface {
	contextExtractor := httpContextValuesExtractor.NewHttpContextValuesExtractor()

	return map[string]port.HandlerInterface{
		CreateFollowUserConst: followHandler.NewCreateManageFollowUserHandler(
			connection,
			uuid,
			contextExtractor,
		),
		// === Register handler marker ===
	}
}
