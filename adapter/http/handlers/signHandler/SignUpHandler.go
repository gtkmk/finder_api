package signHandler

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/userDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	"github.com/gtkmk/finder_api/core/usecase/sharedMethods"
	userUsecase "github.com/gtkmk/finder_api/core/usecase/user"
	"github.com/gtkmk/finder_api/infra/database/repository"
	"github.com/gtkmk/finder_api/infra/envMode"
	"github.com/gtkmk/finder_api/infra/requestEntity/userRequestEntity"
)

type SignUpHandler struct {
	Connection          port.ConnectionInterface
	NotificationService port.NotificationInterface
	PasswordEncryptor   port.EncryptionInterface
	Uuid                port.UuidInterface
	ContextExtractor    port.HttpContextValuesExtractorInterface
	UserDatabase        repositories.UserRepository
	CustomError         port.CustomErrorInterface
	userEventDatabase   repositories.UserEventRepositoryInterface
}

func NewSignUpHandler(
	connection port.ConnectionInterface,
	notificationService port.NotificationInterface,
	passwordEncryptor port.EncryptionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &SignUpHandler{
		Connection:          connection,
		NotificationService: notificationService,
		PasswordEncryptor:   passwordEncryptor,
		Uuid:                uuid,
		ContextExtractor:    contextExtractor,
		CustomError:         customError.NewCustomError(),
	}
}

func (signUpHandler *SignUpHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, signUpHandler.Connection, signUpHandler.Uuid)

	user, err := signUpHandler.defineUser(context)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)

		return
	}

	transaction, err := signUpHandler.Connection.BeginTransaction()

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			signUpHandler.CustomError.ThrowError(err.Error()),
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	signUpHandler.openTableConnection(transaction)

	userEvent := sharedMethods.NewPersistUserEvent(signUpHandler.userEventDatabase, signUpHandler.Uuid, signUpHandler.CustomError)

	ip := context.ClientIP()
	userAgent := context.GetHeader("user-Agent")

	if err := userUsecase.NewCreateUser(
		transaction,
		signUpHandler.UserDatabase,
		user,
		signUpHandler.Uuid,
		signUpHandler.NotificationService,
		signUpHandler.generateUrlResetPassword(user.Id),
		customError.NewCustomError(),
		userEvent,
	).Execute(ip, userAgent); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)

		return
	}

	if err != nil {
		jsonResponse.ThrowError(routesConstants.MessageKeyConst, signUpHandler.CustomError.ThrowError(err.Error()), routesConstants.BadRequestConst)
		return
	}

	jsonResponse.SendJson(routesConstants.MessageKeyConst, "Usuário criado com sucesso", routesConstants.StatusOk)
}

func (signUpHandler *SignUpHandler) openTableConnection(transaction port.ConnectionInterface) {
	signUpHandler.UserDatabase = repository.NewUserDatabase(transaction)
	signUpHandler.userEventDatabase = repository.NewUserEventRepository(transaction)
}

func (signUpHandler *SignUpHandler) defineUser(context *gin.Context) (*userDomain.User, error) {
	userRequest, err := userRequestEntity.SignUpDecodeUserRequest(context.Request)

	if err != nil {
		return nil, signUpHandler.CustomError.ThrowError(err.Error())
	}

	if err := userRequest.Validate(); err != nil {
		return nil, signUpHandler.CustomError.ThrowError(err.Error())
	}

	userEncryptedPwd, err := signUpHandler.PasswordEncryptor.GenerateHashPassword(userRequest.Password)

	if err != nil {
		return nil, signUpHandler.CustomError.ThrowError(err.Error())
	}

	createdAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return nil, signUpHandler.CustomError.ThrowError(err.Error())
	}

	return userDomain.NewUser(
		signUpHandler.Uuid.GenerateUuid(),
		userRequest.Name,
		userRequest.UserName,
		userRequest.Email,
		userEncryptedPwd,
		userRequest.Cpf,
		userRequest.CellphoneNumber,
		userDomain.UserStatusLoggedConst,
		true,
		false,
		createdAt,
	), nil
}

func (signUpHandler *SignUpHandler) generateUrlResetPassword(
	userId string,
) string {
	return fmt.Sprintf("%s/signin?user-id=%s", os.Getenv(envMode.FrontUrlConst), userId)
}
