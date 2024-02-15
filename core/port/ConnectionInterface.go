package port

type ConnectionInterface interface {
	Open() error
	Raw(sql string, statement interface{}, values ...any) error
	Rows(sql string, values ...any) ([]map[string]interface{}, error)
	Close() error
	BeginTransaction() (ConnectionInterface, error)
	SavePoint(checkPointName string) error
	RollbackTo(checkPointName string) error
	Commit() error
	Rollback() error
}
