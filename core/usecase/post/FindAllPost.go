package postUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/filterDomain"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	sharedMethodsInterface "github.com/gtkmk/finder_api/core/port/sharedMethods"
)

const (
	maxPageItensConst = 15
)

type FindPostAll struct {
	PostDatabase         repositories.PostRepositoryInterface
	PostsFilters         *filterDomain.PostFilter
	CalculateQueryOffset sharedMethodsInterface.QueryOffsetInterface
	CustomError          port.CustomErrorInterface
}

func NewFindPostAll(
	postDatabase repositories.PostRepositoryInterface,
	postsFilters *filterDomain.PostFilter,
	calculateQueryOffset sharedMethodsInterface.QueryOffsetInterface,
	customError port.CustomErrorInterface,
) *FindPostAll {
	return &FindPostAll{
		PostDatabase:         postDatabase,
		PostsFilters:         postsFilters,
		CalculateQueryOffset: calculateQueryOffset,
		CustomError:          customError,
	}
}

func (findPostAll *FindPostAll) Execute() ([]map[string]interface{}, error) {
	var err error
	findPostAll.PostsFilters.OffSet, err = findPostAll.CalculateQueryOffset.CalculateQueryOffset(
		helper.ConvertToString(findPostAll.PostsFilters.Limit),
		helper.ConvertToString(*findPostAll.PostsFilters.Page),
	)

	if err != nil {
		return nil, err
	}

	return findPostAll.PostDatabase.FindAllPosts(
		findPostAll.PostsFilters,
	)
}
