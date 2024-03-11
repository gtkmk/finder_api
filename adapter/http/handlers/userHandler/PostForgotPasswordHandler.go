package userHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/success"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	userUsecase "github.com/gtkmk/finder_api/core/usecase/user"
	"github.com/gtkmk/finder_api/infra/database/repository"
	"github.com/gtkmk/finder_api/infra/requestEntity/userRequestEntity"
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
