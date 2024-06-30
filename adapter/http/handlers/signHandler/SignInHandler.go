package signHandler

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/credentialsDomain"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	"github.com/gtkmk/finder_api/core/usecase/sharedMethods"
	userUsecase "github.com/gtkmk/finder_api/core/usecase/user"
	"github.com/gtkmk/finder_api/infra/database/repository"
	"github.com/gtkmk/finder_api/infra/envMode"
	"github.com/gtkmk/finder_api/infra/requestEntity/userRequestEntity"
)

type SignInHandler struct {
	connection        port.ConnectionInterface
	userDatabase      repositories.UserRepository
	userEventDatabase repositories.UserEventRepositoryInterface
	passwordEncryptor port.EncryptionInterface
	uuid              port.UuidInterface
	port.CustomErrorInterface
}

func NewSignInHandler(
	connection port.ConnectionInterface,
	passwordEncryptor port.EncryptionInterface,
	uuid port.UuidInterface,
) port.HandlerInterface {
	return &SignInHandler{
		connection:           connection,
		passwordEncryptor:    passwordEncryptor,
		uuid:                 uuid,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (signInHandler *SignInHandler) Handle(context *gin.Context) {
	fmt.Println("(((((((((((())))))))))))")
	jsonResponse := routes.NewJsonResponse(context, signInHandler.connection, signInHandler.uuid)

	credentials, err := signInHandler.defineCredentials(context)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	signInHandler.openTableConnection()

	userEvent := sharedMethods.NewPersistUserEvent(signInHandler.userEventDatabase, signInHandler.uuid, signInHandler.CustomErrorInterface)

	signIn := userUsecase.NewSign(signInHandler.userDatabase, signInHandler.passwordEncryptor, credentials, userEvent)

	ip := context.ClientIP()
	userAgent := context.GetHeader("user-Agent")

	jwtToken, err := signIn.Execute(ip, userAgent)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	signInHandler.writeCookie(context, jwtToken)
	jsonResponse.SendJson(routesConstants.MessageKeyConst, "Logado com sucesso!", routesConstants.StatusOk)
}

func (signInHandler *SignInHandler) openTableConnection() {
	signInHandler.userDatabase = repository.NewUserDatabase(signInHandler.connection)
	signInHandler.userEventDatabase = repository.NewUserEventRepository(signInHandler.connection)
}

func (signInHandler *SignInHandler) writeCookie(context *gin.Context, token string) {
	context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	fmt.Println("token: ", token)
	domain := os.Getenv(envMode.ApplicationDomainConst)
	context.SetCookie(
		"token",
		token,
		3600, // one hour
		"/",
		domain,
		true,
		true,
	)
}

func (signInHandler *SignInHandler) defineCredentials(c *gin.Context) (*credentialsDomain.Credential, error) {
	credentialRequest, err := userRequestEntity.SignInDecodeUserRequest(c.Request)

	if err != nil {
		return nil, signInHandler.ThrowError(err.Error())
	}

	credentials := credentialsDomain.NewCredentials(credentialRequest.Email, credentialRequest.Password)

	return credentials, nil
}
