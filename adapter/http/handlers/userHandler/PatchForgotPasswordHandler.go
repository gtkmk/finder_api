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

type PatchForgotPasswordHandler struct {
	connection        port.ConnectionInterface
	uuid              port.UuidInterface
	userDatabase      repositories.UserRepository
	passwordEncryptor port.EncryptionInterface
	port.CustomErrorInterface
}

func NewPatchForgotPasswordHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	passwordEncryptor port.EncryptionInterface,
) port.HandlerInterface {
	return &PatchForgotPasswordHandler{
		connection:           connection,
		uuid:                 uuid,
		passwordEncryptor:    passwordEncryptor,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (patchForgotPasswordHandler *PatchForgotPasswordHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, patchForgotPasswordHandler.connection, patchForgotPasswordHandler.uuid)

	patchForgotPasswordHandler.openTableConnection()

	passwordRequest, reqError := userRequestEntity.ForgotUserPasswordRequest(context.Request)

	if reqError != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			patchForgotPasswordHandler.ThrowError(reqError.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	if err := passwordRequest.ThrowsErrorIfPasswordsDoesNotMatch(); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			patchForgotPasswordHandler.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	if err := userUsecase.NewForgotPassword(
		patchForgotPasswordHandler.userDatabase,
		patchForgotPasswordHandler.passwordEncryptor,
	).Execute(passwordRequest.Id, passwordRequest.Password); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson(routesConstants.MessageKeyConst, success.SuccessfullyUpdatedPasswordConst, routesConstants.StatusOk)
}

func (patchForgotPasswordHandler *PatchForgotPasswordHandler) openTableConnection() {
	patchForgotPasswordHandler.userDatabase = repository.NewUserDatabase(
		patchForgotPasswordHandler.connection,
	)
}
