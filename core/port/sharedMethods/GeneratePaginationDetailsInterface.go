package sharedMethods

type GeneratePaginationDetailsInterface interface {
	GeneratePaginationDetails(totalItens, limit, page int64, data []map[string]interface{}) (map[string]interface{}, error)
	MapDBPostToPaginationDetails(dbProposal map[string]interface{}) (map[string]interface{}, error)
}
