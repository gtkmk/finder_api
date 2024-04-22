package sharedMethods

import (
	"fmt"
	"time"

	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/sharedMethods"
)

const (
	PaginationFieldLimitConst = "limite"
)

type GeneratePaginationDetails struct {
	customError port.CustomErrorInterface
}

func NewGeneratePaginationDetails(
	customError port.CustomErrorInterface,
) sharedMethods.GeneratePaginationDetailsInterface {
	return &GeneratePaginationDetails{
		customError: customError,
	}
}

func (generatePaginationDetails *GeneratePaginationDetails) GeneratePaginationDetails(
	totalItems,
	limit,
	page int64,
	data []map[string]interface{},
) (map[string]interface{}, error) {
	if limit == 0 {
		return nil, generatePaginationDetails.customError.ThrowError(
			helper.FieldIsMandatoryAndMustToBeGreaterThanZeroConst,
			PaginationFieldLimitConst,
		)
	}

	if totalItems == 0 {
		return generatePaginationDetails.buildEmptyPaginationDetails(limit, page, data), nil
	}

	totalPages := totalItems / limit
	if totalItems%limit != 0 {
		totalPages++
	}

	if page < 1 || page > totalPages {
		return nil, generatePaginationDetails.customError.ThrowError(
			fmt.Sprintf(helper.InvalidPageNumberErrorConst),
		)
	}

	hasPrevPage := page > 1
	hasNextPage := page < totalPages

	var prevPage, nextPage interface{}
	if hasPrevPage {
		prevPage = page - 1
	} else {
		prevPage = nil
	}
	if hasNextPage {
		nextPage = page + 1
	}

	paginationDetails := map[string]interface{}{
		"totalItems":  totalItems,
		"limit":       limit,
		"totalPages":  totalPages,
		"page":        page,
		"hasPrevPage": hasPrevPage,
		"hasNextPage": hasNextPage,
		"prevPage":    prevPage,
		"nextPage":    nextPage,
		"data":        data,
	}

	return paginationDetails, nil
}

func (generatePaginationDetails *GeneratePaginationDetails) buildEmptyPaginationDetails(limit, page int64, data []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"totalItems":  0,
		"limit":       limit,
		"totalPages":  0,
		"page":        page,
		"hasPrevPage": false,
		"hasNextPage": false,
		"prevPage":    nil,
		"nextPage":    nil,
		"data":        data,
	}
}

func (generatePaginationDetails *GeneratePaginationDetails) MapDBPostToPaginationDetails(dbPost map[string]interface{}) (map[string]interface{}, error) {
	postDate, dateErr := datetimeDomain.FormatDateAsTimeReverted(dbPost["created_at"].(string))

	if dateErr != nil {
		return nil, generatePaginationDetails.customError.ThrowError(dateErr.Error())
	}

	return map[string]interface{}{
		"post_id":            dbPost["post_id"].(string),
		"post_author":        dbPost["post_author"].(string),
		"post_author_avatar": dbPost["post_author_avatar"].(string),
		"text":               dbPost["text"].(string),
		"created_at":         postDate,
		"post_location":      dbPost["post_location"],
		"post_media":         dbPost["post_media"].(string),
		"likes":              dbPost["likes"].(int64),
		"shares":             dbPost["shares"].(int64),
		"comments":           dbPost["comments"].(int64),
		"post_category":      dbPost["post_category"].(string),
		"post_reward":        dbPost["post_reward"].(int64) > 0,
		"post_lostFound":     dbPost["post_lost_found"].(string),
	}, nil
}

func (generatePaginationDetails *GeneratePaginationDetails) MapDBCommentsToPaginationDetails(dbComments map[string]interface{}) (map[string]interface{}, error) {
	creationDate, updateDate, dateErr := generatePaginationDetails.formatDates(
		dbComments["created_at"].(string),
		dbComments["updated_at"],
	)

	if dateErr != nil {
		return nil, generatePaginationDetails.customError.ThrowError(dateErr.Error())
	}

	return map[string]interface{}{
		"comment_id":            dbComments["comment_id"].(string),
		"author_name":           dbComments["author_name"].(string),
		"author_username":       dbComments["author_username"].(string),
		"comment_author_avatar": dbComments["comment_author_avatar"].(string),
		"text":                  dbComments["text"].(string),
		"created_at":            creationDate,
		"updated_at":            updateDate,
		"likes":                 dbComments["likes"].(int64),
	}, nil
}

func (generatePaginationDetails *GeneratePaginationDetails) formatDates(createdDateStr string, updatedDateInterface interface{}) (time.Time, time.Time, error) {
	creationDate, creationErr := datetimeDomain.FormatDateAsTimeReverted(createdDateStr)
	if creationErr != nil {
		return time.Time{}, time.Time{}, creationErr
	}

	if updatedDateInterface != nil {
		updateDate, updateErr := datetimeDomain.FormatDateAsTimeReverted(updatedDateInterface.(string))
		if updateErr != nil {
			return time.Time{}, time.Time{}, updateErr
		}

		return creationDate, updateDate, nil
	}

	return creationDate, time.Time{}, nil
}
