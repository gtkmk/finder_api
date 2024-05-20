package postHandler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	handlerSharedMethods "github.com/gtkmk/finder_api/adapter/http/handlers/sharedMethods"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/filterDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	postUsecase "github.com/gtkmk/finder_api/core/usecase/post"
	"github.com/gtkmk/finder_api/core/usecase/sharedMethods"
	"github.com/gtkmk/finder_api/infra/database/repository"
	"github.com/gtkmk/finder_api/infra/requestEntity/postRequestEntity"
)

type FindPostAllHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	postDatabase     repositories.PostRepositoryInterface
	ContextExtractor port.HttpContextValuesExtractorInterface
	customError      port.CustomErrorInterface
}

func NewFindPostAllHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &FindPostAllHandler{
		connection:       connection,
		uuid:             uuid,
		ContextExtractor: contextExtractor,
		customError:      customError.NewCustomError(),
	}
}

func (findPostAllHandler *FindPostAllHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, findPostAllHandler.connection, findPostAllHandler.uuid)

	loggedUserId, extractErr := findPostAllHandler.ContextExtractor.Extract(context)
	if extractErr != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findPostAllHandler.customError.ThrowError(extractErr.Error()),
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	postsFilters, err := findPostAllHandler.definePostsFilter(context, loggedUserId)
	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findPostAllHandler.customError.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	isOwnProfile := context.Query("is_own_profile")

	fmt.Println("kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk")
	fmt.Println(isOwnProfile)
	fmt.Println("*********************************")
	if isOwnProfile != "" {
		postsFilters.UserId = &loggedUserId
	}

	calculateQueryOffsetSharedMethod := sharedMethods.NewCalculateQueryOffset()

	findPostAllHandler.openTableConnection()

	posts, err := postUsecase.NewFindPostAll(
		findPostAllHandler.postDatabase,
		postsFilters,
		calculateQueryOffsetSharedMethod,
		findPostAllHandler.customError,
	).Execute()

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findPostAllHandler.customError.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	responsePosts, transformErr := findPostAllHandler.transformPostsIntoArrayOfPostDomain(
		posts,
		postsFilters,
	)

	if transformErr != nil {
		jsonResponse.ThrowError(routesConstants.MessageKeyConst, transformErr, routesConstants.BadRequestConst)
		return
	}

	jsonResponse.SendJson(routesConstants.DataKeyConst, map[string]interface{}{"posts": responsePosts}, routesConstants.StatusOk)
}

func (findPostAllHandler *FindPostAllHandler) definePostsFilter(
	context *gin.Context,
	loggedUserId string,
) (
	*filterDomain.PostFilter,
	error,
) {
	checkForSqlInjectionSharedMethod := handlerSharedMethods.NewCheckForSqlInjection(findPostAllHandler.customError)

	decodedRequest, err := postRequestEntity.NewListPostsRequest(
		context,
		findPostAllHandler.uuid,
		checkForSqlInjectionSharedMethod,
		loggedUserId,
	)
	if err != nil {
		return nil, findPostAllHandler.customError.ThrowError(err.Error())
	}

	if err := decodedRequest.ValidatePostsFilterFields(context); err != nil {
		return nil, findPostAllHandler.customError.ThrowError(err.Error())
	}

	return decodedRequest.ConvertProposalFiltersIntoFilterDomain(), nil
}

func (findPostAllHandler *FindPostAllHandler) openTableConnection() {
	findPostAllHandler.postDatabase = repository.NewPostDatabase(findPostAllHandler.connection)
}

func (findPostAllHandler *FindPostAllHandler) transformPostsIntoArrayOfPostDomain(
	dbProposals []map[string]interface{},
	filter *filterDomain.PostFilter,
) (map[string]interface{}, error) {
	var proposals []map[string]interface{}

	generatePaginationDetails := handlerSharedMethods.NewGeneratePaginationDetails(findPostAllHandler.customError)

	if len(dbProposals) == 0 {
		return generatePaginationDetails.GeneratePaginationDetails(0, filter.Limit, *filter.Page, []map[string]interface{}{})
	}

	proposals = make([]map[string]interface{}, len(dbProposals))
	for i, value := range dbProposals {
		proposal, err := generatePaginationDetails.MapDBPostToPaginationDetails(value)
		if err != nil {
			return nil, err
		}

		proposals[i] = proposal
	}

	totalProposals := dbProposals[0]["total_records"].(int64)
	limit := filter.Limit
	page := *filter.Page

	return generatePaginationDetails.GeneratePaginationDetails(totalProposals, limit, page, proposals)
}
