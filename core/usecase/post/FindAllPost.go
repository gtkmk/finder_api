package postUsecase

import (
	"strconv"

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
		PostDatabase: postDatabase,
		CustomError:  customError,
	}
}

func (findPostAll *FindPostAll) Execute(postsFilters *filterDomain.PostFilter) ([]map[string]interface{}, error) {
	findPostAll.PostsFilters = postsFilters

	var err error
	findPostAll.PostsFilters.OffSet, err = findPostAll.calculateQueryOffset(
		helper.ConvertToString(postsFilters.Limit),
		helper.ConvertToString(*postsFilters.Page),
	)

	if err != nil {
		return nil, err
	}

	return findPostAll.PostDatabase.FindAllPosts(
		findPostAll.PostsFilters,
	)
}

func (findPostAll *FindPostAll) calculateQueryOffset(limit string, page string) (*int64, error) {
	if page == "" {
		return nil, nil
	}

	intLimit, err := strconv.Atoi(limit)

	if err != nil {
		return nil, err
	}

	intPage, _ := strconv.Atoi(page)

	if err != nil {
		return nil, err
	}

	offset := int64(intLimit * (intPage - 1))
	return &offset, nil
}
