package paginationDomain

import (
	"github.com/gtkmk/finder_api/core/domain/helper"
	
)

const (
	RecordsPerPageConst = 1000
	MinPageCountConst   = 1
	MaxPageCountConst   = 10
)

type Pagination struct {
	Limit  []int
	Offset []int
}

func CalculatePagination(recordCount int) (Pagination, error) {
	if recordCount <= 0 {
		return Pagination{}, helper.ErrorBuilder(helper.NoRecordsFoundConst)
	}

	pageCount := recordCount / RecordsPerPageConst
	if pageCount < MinPageCountConst {
		pageCount = MinPageCountConst
	} else if pageCount > MaxPageCountConst {
		pageCount = MaxPageCountConst
	}

	recordsPerPage := recordCount / pageCount

	pagination := Pagination{
		Limit:  make([]int, pageCount),
		Offset: make([]int, pageCount),
	}

	for i := 0; i < pageCount; i++ {
		pagination.Offset[i] = i * recordsPerPage
		pagination.Limit[i] = recordsPerPage
	}

	remaining := recordCount % pageCount
	for i := 0; i < remaining; i++ {
		pagination.Limit[i]++
	}

	return pagination, nil
}
