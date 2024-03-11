package userHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper/constants"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	userUsecase "github.com/gtkmk/finder_api/core/usecase/user"
	"github.com/gtkmk/finder_api/infra/database/repository"
	"github.com/gtkmk/finder_api/infra/requestEntity/userRequestEntity"
)

type ResetForgotPasswordHandler struct {
	connection        port.ConnectionInterface
	uuid              port.UuidInterface
	userDatabase      repositories.UserRepository
	passwordEncryptor port.EncryptionInterface
	contextExtractor  port.HttpContextValuesExtractorInterface
	port.CustomErrorInterface
}

func NewPatchResetPasswordHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	passwordEncryptor port.EncryptionInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &ResetForgotPasswordHandler{
		connection:           connection,
		uuid:                 uuid,
		passwordEncryptor:    passwordEncryptor,
		contextExtractor:     contextExtractor,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (resetForgotPasswordHandler *ResetForgotPasswordHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, resetForgotPasswordHandler.connection, resetForgotPasswordHandler.uuid)

	loggedUserId, _, extractError := resetForgotPasswordHandler.contextExtractor.Extract(context)

	if extractError != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			resetForgotPasswordHandler.ThrowError(extractError.Error()),
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	resetForgotPasswordHandler.openTableConnection()

	passwordRequest, reqError := userRequestEntity.ResetUserPasswordRequest(context.Request)

	if reqError != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			resetForgotPasswordHandler.ThrowError(reqError.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	if err := passwordRequest.Validate(); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			resetForgotPasswordHandler.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	if passwordRequest.Id != loggedUserId {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			resetForgotPasswordHandler.ThrowError(constants.WithoutPermissionConst),
			routesConstants.ForbiddenRequestConst,
		)
		return
	}

	if err := userUsecase.NewResetPassword(
		resetForgotPasswordHandler.userDatabase,
		resetForgotPasswordHandler.passwordEncryptor,
	).Execute(
		passwordRequest.Id,
		passwordRequest.OldPassword,
		passwordRequest.Password,
	); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson(routesConstants.MessageKeyConst, "Senha atualizada com sucesso!", routesConstants.StatusOk)
}

func (resetForgotPasswordHandler *ResetForgotPasswordHandler) openTableConnection() {
	resetForgotPasswordHandler.userDatabase = repository.NewUserDatabase(
		resetForgotPasswordHandler.connection,
	)
}
