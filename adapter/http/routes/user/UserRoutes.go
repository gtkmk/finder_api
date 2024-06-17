package user

import (
	"os"

	"github.com/gtkmk/finder_api/adapter/http/handlers/signHandler"
	"github.com/gtkmk/finder_api/adapter/http/handlers/userHandler"
	"github.com/gtkmk/finder_api/adapter/http/middleware"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/jwtAuth"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/envMode"
	"github.com/gtkmk/finder_api/infra/httpContextValuesExtractor"

	"github.com/gin-gonic/gin"
)

const (
	GetLoggedUserKeyConst           string = "getLoggedUser"
	GetUserKeyConst                 string = "getUser"
	GetUserListByNameConst          string = "getUserListByName"
	PatchUserFirstAccessKeyConst    string = "patchUserFirstAccess"
	PatchUserForgotPasswordKeyConst string = "patchUserForgotPassword"
	PatchResetUserPasswordKeyConst  string = "patchResetUserPassword"
	PatchUserHandlerKeyConst        string = "patchUserHandler"
	SignInHandlerKeyConst           string = "signInHandler"
	SignOutHandlerKeyConst          string = "signOutHandler"
	SignUpHandlerKeyConst           string = "signUpHandler"
	FindUserDetailsConst            string = "findUserDetails"
	// === Route constants marker ===
	UpdateUserUpdateUserInfoConst string = "UpdateUser"
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
		routesConstants.PostSignInRouteConst,
		userRoutes.userHandlers[SignInHandlerKeyConst].Handle,
	)

	userRoutes.POST(
		routesConstants.PostSignOutRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[SignOutHandlerKeyConst].Handle,
	)

	userRoutes.POST(
		routesConstants.PostSignUpRouteConst,
		userRoutes.userHandlers[SignUpHandlerKeyConst].Handle,
	)

	userRoutes.GET(
		routesConstants.GetLoggedUserRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[GetLoggedUserKeyConst].Handle,
	)

	userRoutes.GET(
		routesConstants.GetUserDetailsRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[FindUserDetailsConst].Handle,
	)

	userRoutes.GET(
		routesConstants.GetUsersListByNameRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[GetUserListByNameConst].Handle,
	)

	//Erro em algum lugar por aqui:

	//userRoutes.GET(
	//	routesConstants.GetUserRouteConst,
	//	userRoutes.jwt.IsAuthorizedMiddleware(),
	//	userRoutes.userHandlers[GetUserKeyConst].Handle,
	//)
	//
	//userRoutes.PATCH(
	//	routesConstants.PatchResetPasswordRouteConst,
	//	userRoutes.jwt.IsAuthorizedMiddleware(),
	//	userRoutes.userHandlers[PatchResetUserPasswordKeyConst].Handle,
	//)
	//
	//userRoutes.PATCH(
	//	routesConstants.PatchEditUserRouteConst,
	//	userRoutes.jwt.IsAuthorizedMiddleware(),
	//	userRoutes.userHandlers[PatchUserHandlerKeyConst].Handle,
	//)
	//
	//userRoutes.POST(
	//	routesConstants.PostSignInRouteConst,
	//	userRoutes.userHandlers[SignInHandlerKeyConst].Handle,
	//)
	//
	//userRoutes.PATCH(
	//	routesConstants.PatchFirstAccessRouteConst,
	//	userRoutes.userHandlers[PatchUserFirstAccessKeyConst].Handle,
	//)
	//
	//userRoutes.PATCH(
	//	routesConstants.PatchForgotPasswordRouteConst,
	//	userRoutes.userHandlers[PatchUserForgotPasswordKeyConst].Handle,
	//)

	// === Register route marker ===
	userRoutes.PATCH(
		routesConstants.PatchUserInfoRouteConst,
		userRoutes.jwt.IsAuthorizedMiddleware(),
		userRoutes.userHandlers[UpdateUserUpdateUserInfoConst].Handle,
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
		PatchUserFirstAccessKeyConst: userHandler.NewPatchFirstAccessHandler(
			connection,
			uuid,
			passwordEncryption,
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
		GetUserKeyConst: userHandler.NewFindUserHandler(
			connection,
			uuid,
		),
		FindUserDetailsConst: userHandler.NewFindUserUserDetailsHandler(
			connection,
			uuid,
			contextExtractor,
		),
		GetUserListByNameConst: userHandler.NewFindUsersListByNameHandler(
			connection,
			uuid,
			contextExtractor,
		),
		// === Register handler marker ===
		UpdateUserUpdateUserInfoConst: userHandler.NewUpdateUserUpdateUserInfoHandler(connection, uuid, contextExtractor),
	}
}
