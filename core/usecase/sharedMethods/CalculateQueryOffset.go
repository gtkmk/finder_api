package sharedMethods

import (
	"strconv"
)

type CalculateQueryOffset struct{}

func NewCalculateQueryOffset() *CalculateQueryOffset {
	return &CalculateQueryOffset{}
}

func (calculateQueryOffset *CalculateQueryOffset) CalculateQueryOffset(limit string, page string) (*int64, error) {
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
