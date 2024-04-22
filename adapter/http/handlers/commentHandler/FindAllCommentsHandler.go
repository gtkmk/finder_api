package commentHandler

import (
	"github.com/gin-gonic/gin"

	handlerSharedMethods "github.com/gtkmk/finder_api/adapter/http/handlers/sharedMethods"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	commentUsecase "github.com/gtkmk/finder_api/core/usecase/comment"
	"github.com/gtkmk/finder_api/core/usecase/sharedMethods"
	"github.com/gtkmk/finder_api/infra/database/repository"
	"github.com/gtkmk/finder_api/infra/requestEntity/commentRequestEntity"
)

const (
	maxPageItensConst = 5
)

type FindAllCommentsHandler struct {
	Connection      port.ConnectionInterface
	Uuid            port.UuidInterface
	CommentDatabase repositories.CommentRepository
	PostDatabase    repositories.PostRepositoryInterface
	CustomError     port.CustomErrorInterface
}

func NewFindAllCommentsHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
) port.HandlerInterface {
	return &FindAllCommentsHandler{
		Connection:  connection,
		Uuid:        uuid,
		CustomError: customError.NewCustomError(),
	}
}

func (findCommentFindAllHandler *FindAllCommentsHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, findCommentFindAllHandler.Connection, findCommentFindAllHandler.Uuid)

	postId, actualPage, err := findCommentFindAllHandler.defineCommentsFilter(context)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findCommentFindAllHandler.CustomError.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	calculateQueryOffsetSharedMethod := sharedMethods.NewCalculateQueryOffset()

	findCommentFindAllHandler.openTableConnection()

	comments, err := commentUsecase.NewFindCommentFindAll(
		postId,
		calculateQueryOffsetSharedMethod,
		findCommentFindAllHandler.CommentDatabase,
		findCommentFindAllHandler.PostDatabase,
		findCommentFindAllHandler.CustomError,
	).Execute(int(actualPage))

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findCommentFindAllHandler.CustomError.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	responseComments, transformErr := findCommentFindAllHandler.transformCommentsIntoPaginatedResponse(
		actualPage,
		comments,
	)

	if transformErr != nil {
		jsonResponse.ThrowError(routesConstants.MessageKeyConst, transformErr, routesConstants.BadRequestConst)
		return
	}

	jsonResponse.SendJson(routesConstants.DataKeyConst, map[string]interface{}{"comments": responseComments}, routesConstants.StatusOk)
}

func (findCommentFindAllHandler *FindAllCommentsHandler) openTableConnection() {
	findCommentFindAllHandler.CommentDatabase = repository.NewCommentDatabase(findCommentFindAllHandler.Connection)
	findCommentFindAllHandler.PostDatabase = repository.NewPostDatabase(findCommentFindAllHandler.Connection)
}

func (findCommentFindAllHandler *FindAllCommentsHandler) defineCommentsFilter(
	context *gin.Context,
) (
	string,
	int64,
	error,
) {
	decodedRequest, err := commentRequestEntity.NewListCommentsRequest(
		context,
		findCommentFindAllHandler.Uuid,
	)
	if err != nil {
		return "", 0, findCommentFindAllHandler.CustomError.ThrowError(err.Error())
	}

	if err := decodedRequest.ValidateCommentsFilterFields(context); err != nil {
		return "", 0, findCommentFindAllHandler.CustomError.ThrowError(err.Error())
	}

	postId, page := decodedRequest.RetrieveCommentsFiltersInfo()
	return postId, page, nil
}

func (findCommentFindAllHandler *FindAllCommentsHandler) transformCommentsIntoPaginatedResponse(
	page int64,
	dbComments []map[string]interface{},
) (map[string]interface{}, error) {
	var comments []map[string]interface{}

	generatePaginationDetails := handlerSharedMethods.NewGeneratePaginationDetails(findCommentFindAllHandler.CustomError)

	if len(dbComments) == 0 {
		return generatePaginationDetails.GeneratePaginationDetails(0, maxPageItensConst, page, []map[string]interface{}{})
	}

	comments = make([]map[string]interface{}, len(dbComments))
	for i, value := range dbComments {
		proposal, err := generatePaginationDetails.MapDBCommentsToPaginationDetails(value)
		if err != nil {
			return nil, err
		}

		comments[i] = proposal
	}

	totalComments := dbComments[0]["total_records"].(int64)

	return generatePaginationDetails.GeneratePaginationDetails(totalComments, maxPageItensConst, page, comments)
}
