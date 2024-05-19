package commentUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	sharedMethodsInterface "github.com/gtkmk/finder_api/core/port/sharedMethods"
)

const (
	maxPageItensConst = 5
)

type FindCommentFindAll struct {
	PostId               string
	CalculateQueryOffset sharedMethodsInterface.QueryOffsetInterface
	CommentDatabase      repositories.CommentRepository
	PostDatabase         repositories.PostRepositoryInterface
	CustomError          port.CustomErrorInterface
}

func NewFindCommentFindAll(
	postId string,
	calculateQueryOffset sharedMethodsInterface.QueryOffsetInterface,
	commentDatabase repositories.CommentRepository,
	postDatabase repositories.PostRepositoryInterface,
	customError port.CustomErrorInterface,
) *FindCommentFindAll {
	return &FindCommentFindAll{
		PostId:               postId,
		CalculateQueryOffset: calculateQueryOffset,
		CommentDatabase:      commentDatabase,
		PostDatabase:         postDatabase,
		CustomError:          customError,
	}
}

func (findCommentFindAll *FindCommentFindAll) Execute(actualPage int, loggedUserId string) ([]map[string]interface{}, error) {
	if err := findCommentFindAll.verifyIfPostExists(); err != nil {
		return nil, err
	}

	offSet, err := findCommentFindAll.CalculateQueryOffset.CalculateQueryOffset(
		helper.ConvertToString(maxPageItensConst),
		helper.ConvertToString(actualPage),
	)

	if err != nil {
		return nil, err
	}

	return findCommentFindAll.CommentDatabase.FindAllComments(
		findCommentFindAll.PostId,
		loggedUserId,
		*offSet,
		maxPageItensConst,
	)
}

func (findCommentFindAll *FindCommentFindAll) verifyIfPostExists() error {
	post, err := findCommentFindAll.PostDatabase.FindPostByID(findCommentFindAll.PostId)

	if err != nil {
		return findCommentFindAll.CustomError.ThrowError(err.Error())
	}

	if post == nil {
		return findCommentFindAll.CustomError.ThrowError(helper.PostNotFoundMessageConst)
	}

	return nil
}
