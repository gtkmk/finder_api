package document

import (
	"os"

	"github.com/gtkmk/finder_api/adapter/http/handlers/documentHandler"
	"github.com/gtkmk/finder_api/adapter/http/middleware"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/jwtAuth"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/envMode"

	"github.com/gin-gonic/gin"
)

const (
	// === Route constants marker ===
	FindDocumentImageBase64Const string = "FindDocument"
)

type DocumentRoutes struct {
	*gin.Engine
	documentHandlers map[string]port.HandlerInterface
	jwt              *middleware.IsAuthorized
}

func NewDocumentRoutes(
	app *gin.Engine,
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	passwordEncryption port.EncryptionInterface,
	notificationService port.NotificationInterface,
) port.RoutesInterface {
	jwt := jwtAuth.NewjwtAuth(os.Getenv(envMode.JwtSecretConst))

	return &DocumentRoutes{
		app,
		createMapOfDocumentHandlers(connection, notificationService, uuid, passwordEncryption),
		middleware.NewIsAuthorized(
			jwt,
			connection,
			uuid,
		),
	}
}

func (documentRoutes *DocumentRoutes) Register() {
	// === Register route marker ===
	documentRoutes.GET(
		routesConstants.FindDocumentBase64RouteConst,
		documentRoutes.jwt.IsAuthorizedMiddleware(),
		documentRoutes.documentHandlers[FindDocumentImageBase64Const].Handle,
	)

}

func createMapOfDocumentHandlers(
	connection port.ConnectionInterface,
	notificationService port.NotificationInterface,
	uuid port.UuidInterface,
	passwordEncryption port.EncryptionInterface,
) map[string]port.HandlerInterface {
	// contextExtractor := httpContextValuesExtractor.NewHttpContextValuesExtractor()

	return map[string]port.HandlerInterface{
		// === Register handler marker ===
		FindDocumentImageBase64Const: documentHandler.NewFindDocumentImageBase64Handler(connection, uuid),
	}
}
