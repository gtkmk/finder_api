package comment

import (
	"os"

	"github.com/gtkmk/finder_api/adapter/http/handlers/commentHandler"
	"github.com/gtkmk/finder_api/adapter/http/middleware"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/jwtAuth"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/envMode"
	"github.com/gtkmk/finder_api/infra/httpContextValuesExtractor"

	"github.com/gin-gonic/gin"
)

const (
	CreateCommentConst string = "CreateComment"
	EditCommentConst   string = "EditComment"
	DeleteCommentConst string = "DeleteComment"
	// === Route constants marker ===
	FindCommentFindAllConst string = "FindComment"
)

type CommentRoutes struct {
	*gin.Engine
	commentHandlers map[string]port.HandlerInterface
	jwt             *middleware.IsAuthorized
}

func NewCommentRoutes(
	app *gin.Engine,
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	passwordEncryption port.EncryptionInterface,
	notificationService port.NotificationInterface,
) port.RoutesInterface {
	jwt := jwtAuth.NewjwtAuth(os.Getenv(envMode.JwtSecretConst))

	return &CommentRoutes{
		app,
		createMapOfCommentHandlers(connection, notificationService, uuid, passwordEncryption),
		middleware.NewIsAuthorized(
			jwt,
			connection,
			uuid,
		),
	}
}

func (commentRoutes *CommentRoutes) Register() {
	commentRoutes.POST(
		routesConstants.PostCreateCommentRouteConst,
		commentRoutes.jwt.IsAuthorizedMiddleware(),
		commentRoutes.commentHandlers[CreateCommentConst].Handle,
	)
	commentRoutes.PATCH(
		routesConstants.PatchEditCommentRouteConst,
		commentRoutes.jwt.IsAuthorizedMiddleware(),
		commentRoutes.commentHandlers[EditCommentConst].Handle,
	)
	commentRoutes.DELETE(
		routesConstants.DeleteCommentRouteConst,
		commentRoutes.jwt.IsAuthorizedMiddleware(),
		commentRoutes.commentHandlers[DeleteCommentConst].Handle,
	)
	commentRoutes.GET(
		routesConstants.FindFindAllCommentsRouteConst,
		commentRoutes.jwt.IsAuthorizedMiddleware(),
		commentRoutes.commentHandlers[FindCommentFindAllConst].Handle,
	)

	// === Register route marker ===
}

func createMapOfCommentHandlers(
	connection port.ConnectionInterface,
	notificationService port.NotificationInterface,
	uuid port.UuidInterface,
	passwordEncryption port.EncryptionInterface,
) map[string]port.HandlerInterface {
	contextExtractor := httpContextValuesExtractor.NewHttpContextValuesExtractor()

	return map[string]port.HandlerInterface{
		CreateCommentConst: commentHandler.NewCreateCommentHandler(
			connection,
			uuid,
			contextExtractor,
		),
		EditCommentConst: commentHandler.NewUpdateCommentHandler(
			connection,
			uuid,
			contextExtractor,
		),
		DeleteCommentConst: commentHandler.NewDeleteCommentHandler(
			connection,
			uuid,
			contextExtractor,
		),
		FindCommentFindAllConst: commentHandler.NewFindAllCommentsHandler(
			connection,
			uuid,
		),
		// === Register handler marker ===
	}
}
