package port

type SpreadsheetFactoryInterface interface {
	Make(string, string) SpreadsheetInterface
}
