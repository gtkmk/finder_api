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
