package userHandler

import (
	"github.com/em-cash/simulador.em.cash/adapter/http/routes"
	"github.com/em-cash/simulador.em.cash/adapter/http/routesConstants"
	"github.com/em-cash/simulador.em.cash/core/domain/customError"
	"github.com/em-cash/simulador.em.cash/core/domain/success"
	"github.com/em-cash/simulador.em.cash/core/port"
	"github.com/em-cash/simulador.em.cash/core/port/repositories"
	userUsecase "github.com/em-cash/simulador.em.cash/core/usecase/user"
	"github.com/em-cash/simulador.em.cash/infra/database/repository"
	"github.com/em-cash/simulador.em.cash/infra/requestEntity/userRequestEntity"
	"github.com/gin-gonic/gin"
)

type PostForgotPasswordHandler struct {
	connection          port.ConnectionInterface
	uuid                port.UuidInterface
	notificationService port.NotificationInterface
	userDatabase        repositories.UserRepository
	port.CustomErrorInterface
}

func NewPostForgotPasswordHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	notificationService port.NotificationInterface,
) port.HandlerInterface {
	return &PostForgotPasswordHandler{
		connection:           connection,
		uuid:                 uuid,
		notificationService:  notificationService,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (postForgotPasswordHandler *PostForgotPasswordHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, postForgotPasswordHandler.connection, postForgotPasswordHandler.uuid)

	postForgotPasswordHandler.openTableConnection()

	forgotPasswordRequest, err := userRequestEntity.ForgotPasswordDecodeRequest(context.Request)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			postForgotPasswordHandler.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	forgotPassword := userUsecase.NewForgotPasswordEmail(
		postForgotPasswordHandler.userDatabase,
		postForgotPasswordHandler.notificationService,
	)

	if err := forgotPassword.Execute(forgotPasswordRequest.Email); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson(routesConstants.MessageKeyConst, success.SuccessfullySendResetPasswordEmailConst, routesConstants.StatusOk)
}

func (postForgotPasswordHandler *PostForgotPasswordHandler) openTableConnection() {
	postForgotPasswordHandler.userDatabase = repository.NewUserDatabase(
		postForgotPasswordHandler.connection,
	)
}
