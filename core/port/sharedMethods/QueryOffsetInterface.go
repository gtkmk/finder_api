package sharedMethods

type QueryOffsetInterface interface {
	CalculateQueryOffset(limit string, page string) (*int64, error)
}
