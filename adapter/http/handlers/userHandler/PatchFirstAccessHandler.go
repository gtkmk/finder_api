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

type PatchFirstAccessHandler struct {
	connection        port.ConnectionInterface
	uuid              port.UuidInterface
	userDatabase      repositories.UserRepository
	passwordEncryptor port.EncryptionInterface
	port.CustomErrorInterface
}

func NewPatchFirstAccessHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	passwordEncryptor port.EncryptionInterface,
) port.HandlerInterface {
	return &PatchFirstAccessHandler{
		connection:           connection,
		uuid:                 uuid,
		passwordEncryptor:    passwordEncryptor,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (patchFirstAccessHandler *PatchFirstAccessHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, patchFirstAccessHandler.connection, patchFirstAccessHandler.uuid)

	patchFirstAccessHandler.openTableConnection()

	passwordRequest, reqError := userRequestEntity.ForgotUserPasswordRequest(context.Request)

	if reqError != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			patchFirstAccessHandler.ThrowError(reqError.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	if err := passwordRequest.ThrowsErrorIfPasswordsDoesNotMatch(); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			patchFirstAccessHandler.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	if err := userUsecase.NewUserFirstAccess(
		patchFirstAccessHandler.userDatabase,
		patchFirstAccessHandler.passwordEncryptor,
	).Execute(passwordRequest.Id, passwordRequest.Password); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson(routesConstants.MessageKeyConst, success.SuccessfullyRegisteredNewPasswordConst, routesConstants.StatusOk)
}

func (patchFirstAccessHandler *PatchFirstAccessHandler) openTableConnection() {
	patchFirstAccessHandler.userDatabase = repository.NewUserDatabase(
		patchFirstAccessHandler.connection,
	)
}
