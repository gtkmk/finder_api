package userHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/success"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	"github.com/gtkmk/finder_api/core/usecase/sharedMethods"
	userUsecase "github.com/gtkmk/finder_api/core/usecase/user"
	"github.com/gtkmk/finder_api/infra/database/repository"
	"github.com/gtkmk/finder_api/infra/requestEntity/userRequestEntity"
)

type UpdateUserUpdateUserInfoHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	contextExtractor port.HttpContextValuesExtractorInterface
	userDatabase     repositories.UserRepository
	customError      port.CustomErrorInterface
}

func NewUpdateUserUpdateUserInfoHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &UpdateUserUpdateUserInfoHandler{
		connection:       connection,
		uuid:             uuid,
		contextExtractor: contextExtractor,
		customError:      customError.NewCustomError(),
	}
}

func (updateUserUpdateUserInfoHandler *UpdateUserUpdateUserInfoHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, updateUserUpdateUserInfoHandler.connection, updateUserUpdateUserInfoHandler.uuid)

	loggedUserId, extractErr := updateUserUpdateUserInfoHandler.contextExtractor.Extract(context)
	if extractErr != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			updateUserUpdateUserInfoHandler.customError.ThrowError(extractErr.Error()),
			routesConstants.InternarServerErrorConst,
		)

		return
	}

	decodedUserName, decodedUserCellphone, err := updateUserUpdateUserInfoHandler.defineUser(context)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)

		return
	}

	transaction, err := updateUserUpdateUserInfoHandler.connection.BeginTransaction()

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			updateUserUpdateUserInfoHandler.customError.ThrowError(err.Error()),
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	updateUserUpdateUserInfoHandler.openTableConnection(transaction)

	rollBackAndReturn := sharedMethods.NewRollBackAndReturnError(transaction)

	if err := userUsecase.NewUpdateUserUpdateUserInfo(
		loggedUserId,
		updateUserUpdateUserInfoHandler.userDatabase,
		transaction,
		rollBackAndReturn,
		updateUserUpdateUserInfoHandler.customError,
	).Execute(decodedUserName, decodedUserCellphone); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson(routesConstants.DataKeyConst, success.SuccessfullyUpdatedUserConst, routesConstants.StatusOk)
}

func (updateUserUpdateUserInfoHandler *UpdateUserUpdateUserInfoHandler) defineUser(context *gin.Context) (
	string,
	string,
	error,
) {
	decodedUser, err := userRequestEntity.NewUpdateUserInfoRequest(context)

	if err != nil {
		return "", "", updateUserUpdateUserInfoHandler.customError.ThrowError(err.Error())
	}

	if err := decodedUser.Validate(); err != nil {
		return "", "", updateUserUpdateUserInfoHandler.customError.ThrowError(err.Error())
	}

	decodedRequestInfo, err := decodedUser.DecodeUpdatedUserInfoRequest(context.Request)

	if err != nil {
		return "", "", updateUserUpdateUserInfoHandler.customError.ThrowError(err.Error())
	}

	return decodedRequestInfo.Name, decodedRequestInfo.CellphoneNumber, nil
}

func (updateUserUpdateUserInfoHandler *UpdateUserUpdateUserInfoHandler) openTableConnection(transaction port.ConnectionInterface) {
	updateUserUpdateUserInfoHandler.userDatabase = repository.NewUserDatabase(transaction)
}
