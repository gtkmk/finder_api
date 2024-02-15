package userHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/userDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	userUsecase "github.com/gtkmk/finder_api/core/usecase/user"
	"github.com/gtkmk/finder_api/infra/database/repository"
)

type FindLoggedUserHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	userDatabase     repositories.UserRepository
	contextExtractor port.HttpContextValuesExtractorInterface
	port.CustomErrorInterface
}

func NewFindLoggedUserHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &FindLoggedUserHandler{
		connection:           connection,
		uuid:                 uuid,
		contextExtractor:     contextExtractor,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (findLoggedUserRoute *FindLoggedUserHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, findLoggedUserRoute.connection, findLoggedUserRoute.uuid)

	userId, _, extractError := findLoggedUserRoute.contextExtractor.Extract(context)

	if extractError != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findLoggedUserRoute.ThrowError(extractError.Error()),
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	findLoggedUserRoute.openTableConnection()

	getUser := userUsecase.NewGetUser(findLoggedUserRoute.userDatabase)

	user, err := getUser.Execute(userId, true)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	userReturn := findLoggedUserRoute.defineGetLoggedUserReturn(user)

	jsonResponse.SendJson("user", userReturn, routesConstants.StatusOk)
}

func (findLoggedUserRoute *FindLoggedUserHandler) defineGetLoggedUserReturn(user *userDomain.User) FindUserReturn {
	var userReturn FindUserReturn

	userReturn.Id = user.Id
	userReturn.Email = user.Email
	userReturn.Name = user.Name
	userReturn.IsActive = user.IsActive

	return userReturn
}

func (findLoggedUserRoute *FindLoggedUserHandler) openTableConnection() {
	findLoggedUserRoute.userDatabase = repository.NewUserDatabase(findLoggedUserRoute.connection)
}
