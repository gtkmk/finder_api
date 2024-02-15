package spreadsheet

import (
	"github.com/gtkmk/finder_api/core/domain/documentDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/spreadsheet/csvSpreadsheet"
	"github.com/gtkmk/finder_api/infra/spreadsheet/xlsxSpreadsheet"
)

type SpreadsheetFactory struct{}

func NewSpreadsheetFactory() port.SpreadsheetFactoryInterface {
	return &SpreadsheetFactory{}
}

func (spreadsheetFactory *SpreadsheetFactory) Make(name string, spreadsheetType string) port.SpreadsheetInterface {
	switch spreadsheetType {
	case documentDomain.SpreadsheetTypeCsvConst:
		return csvSpreadsheet.NewCsvSpreadsheet(name)
	case documentDomain.SpreadsheetTypeXlsxConst:
		return xlsxSpreadsheet.NewXlsxSpreadsheet(name)
	default:
		return nil
	}
}
