package user

import (
	"fmt"
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
	GetUserListKeyConst             string = "getUserList"
	PatchUserFirstAccessKeyConst    string = "patchUserFirstAccess"
	PatchUserForgotPasswordKeyConst string = "patchUserForgotPassword"
	PatchResetUserPasswordKeyConst  string = "patchResetUserPassword"
	PatchUserHandlerKeyConst        string = "patchUserHandler"
	PostForgotUserPasswordKeyConst  string = "postForgotUserPassword"
	SignInHandlerKeyConst           string = "signInHandler"
	SignOutHandlerKeyConst          string = "signOutHandler"
	SignUpHandlerKeyConst           string = "signUpHandler"
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

	//Erro em algum lugar por aqui:

	//userRoutes.GET(
	//	routesConstants.GetUserRouteConst,
	//	userRoutes.jwt.IsAuthorizedMiddleware(),
	//	userRoutes.userHandlers[GetUserKeyConst].Handle,
	//)
	//
	//userRoutes.GET(
	//	routesConstants.GetUsersListRouteConst,
	//	userRoutes.jwt.IsAuthorizedMiddleware(),
	//	userRoutes.userHandlers[GetUserListKeyConst].Handle,
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
	//
	//userRoutes.POST(
	//	routesConstants.PostForgotPasswordRouteConst,
	//	userRoutes.userHandlers[PostForgotUserPasswordKeyConst].Handle,
	//)

}

func createMapOfUserHandlers(
	connection port.ConnectionInterface,
	notificationService port.NotificationInterface,
	uuid port.UuidInterface,
	passwordEncryption port.EncryptionInterface,
) map[string]port.HandlerInterface {
	contextExtractor := httpContextValuesExtractor.NewHttpContextValuesExtractor()
	fmt.Println("))))))))))))))))))))))))))))))))))))))))))))))))")
	fmt.Println(connection)
	fmt.Println(notificationService)
	fmt.Println(uuid.GenerateUuid())
	fmt.Println(passwordEncryption.GenerateRandomPassword())
	fmt.Println(signHandler.NewSignInHandler(connection, passwordEncryption, uuid))

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
		PostForgotUserPasswordKeyConst: userHandler.NewPostForgotPasswordHandler(
			connection,
			uuid,
			notificationService,
		),
	}
}
