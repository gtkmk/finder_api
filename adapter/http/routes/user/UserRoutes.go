package user

import (
	"github.com/gtkmk/finder_api/adapter/http/handlers/signHandler"
	"github.com/gtkmk/finder_api/adapter/http/handlers/userHandler"
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
	GetLoggedUserKeyConst                   string = "getLoggedUser"
	GetManagersListKeyConst                 string = "getManagersList"
	GetUserKeyConst                         string = "getUser"
	GetUserListKeyConst                     string = "getUserList"
	GetUserUnitiesConst                     string = "getUserUnities"
	PatchUserFirstAccessKeyConst            string = "patchUserFirstAccess"
	PatchUserForgotPasswordKeyConst         string = "patchUserForgotPassword"
	PatchResetUserPasswordKeyConst          string = "patchResetUserPassword"
	PatchUserPermissionGroupKeyConst        string = "patchUserPermissionGroup"
	PatchUserHandlerKeyConst                string = "patchUserHandler"
	PostForgotUserPasswordKeyConst          string = "postForgotUserPassword"
	PostUserProductsHandlerKeyConst         string = "postUserProductsHandler"
	SignInHandlerKeyConst                   string = "signInHandler"
	SignOutHandlerKeyConst                  string = "signOutHandler"
	SignUpHandlerKeyConst                   string = "signUpHandler"
	GetManagersAndConsultantsListKeyConst   string = "getManagersAndConsultantsList"
	PostUserMassRegistrationHandlerKeyConst string = "postUserMassRegistrationHandler"
	GetExportUsersHandlerKeyConst           string = "getExportUsersHandlerKeyConst"
	GetUserLeaderOptionsConst               string = "getUserLeaderOptionsConst"
)

type UserRoutes struct {
	*gin.Engine
	userHandlers map[string]port.HandlerInterface
	jwt          *middleware.IsAuthorized
}

func NewUserRoutes(
	app *gin.Engine,
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	passwordEncryption port.EncryptionInterface,
	notificationService port.NotificationInterface,
) port.RoutesInterface {
	jwt := jwtAuth.NewjwtAuth(os.Getenv(envMode.JwtSecretConst))

	return &UserRoutes{
		app,
		createMapOfUserHandlers(connection, notificationService, uuid, passwordEncryption),
		middleware.NewIsAuthorized(
			jwt,
			connection,
			uuid,
		),
	}
}

func (userRoutes *UserRoutes) Register() {
	userRoutes.POST(
		routesConstants.PostSignOutRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[SignOutHandlerKeyConst].Handle,
	)

	userRoutes.POST(
		routesConstants.PostSignUpRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[SignUpHandlerKeyConst].Handle,
	)

	userRoutes.GET(
		routesConstants.GetLoggedUserRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[GetLoggedUserKeyConst].Handle,
	)

	userRoutes.GET(
		routesConstants.GetManagersListRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[GetManagersListKeyConst].Handle,
	)

	userRoutes.GET(
		routesConstants.GetUserRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[GetUserKeyConst].Handle,
	)

	userRoutes.GET(
		routesConstants.GetUsersListRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[GetUserListKeyConst].Handle,
	)

	userRoutes.PATCH(
		routesConstants.PatchResetPasswordRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[PatchResetUserPasswordKeyConst].Handle,
	)

	userRoutes.PATCH(
		routesConstants.PatchUserPermissionGroupRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[PatchUserPermissionGroupKeyConst].Handle,
	)

	userRoutes.PATCH(
		routesConstants.PatchEditUserRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[PatchUserHandlerKeyConst].Handle,
	)

	userRoutes.POST(
		routesConstants.PostUserProductRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[PostUserProductsHandlerKeyConst].Handle,
	)

	userRoutes.POST(
		routesConstants.PostSignInRouteConst,
		userRoutes.userHandlers[SignInHandlerKeyConst].Handle,
	)

	userRoutes.PATCH(
		routesConstants.PatchFirstAccessRouteConst,
		userRoutes.userHandlers[PatchUserFirstAccessKeyConst].Handle,
	)

	userRoutes.PATCH(
		routesConstants.PatchForgotPasswordRouteConst,
		userRoutes.userHandlers[PatchUserForgotPasswordKeyConst].Handle,
	)

	userRoutes.POST(
		routesConstants.PostForgotPasswordRouteConst,
		userRoutes.userHandlers[PostForgotUserPasswordKeyConst].Handle,
	)

	userRoutes.POST(
		routesConstants.PostMassRegistrationRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[PostUserMassRegistrationHandlerKeyConst].Handle,
	)

	userRoutes.GET(
		routesConstants.GetManagerAndConsultantsListConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[GetManagersAndConsultantsListKeyConst].Handle,
	)

	userRoutes.GET(
		routesConstants.GetUserUnitiesRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[GetUserUnitiesConst].Handle,
	)

	userRoutes.GET(
		routesConstants.GetExportUsersRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[GetExportUsersHandlerKeyConst].Handle,
	)

	userRoutes.GET(
		routesConstants.GetUserLeaderOptionsConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[GetUserLeaderOptionsConst].Handle,
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
		SignInHandlerKeyConst: signHandler.NewSignInHandler(
			connection,
			passwordEncryption,
			uuid,
		),
		SignOutHandlerKeyConst: signHandler.NewSignOutHandler(
			connection,
			uuid,
		),
		SignUpHandlerKeyConst: signHandler.NewSignUpHandler(
			connection,
			notificationService,
			passwordEncryption,
			uuid,
			contextExtractor,
		),
		GetLoggedUserKeyConst: userHandler.NewFindLoggedUserHandler(
			connection,
			uuid,
			contextExtractor,
		),
		PatchUserForgotPasswordKeyConst: userHandler.NewPatchForgotPasswordHandler(
			connection,
			uuid,
			passwordEncryption,
		),
	}
}
