package port

import (
	"mime/multipart"
)

type SpreadsheetInterface interface {
	GenerateSpreadsheetFromData(data []map[string]interface{}) (*multipart.FileHeader, []byte, error)
	ReadSpreadsheet(filePath string) ([]map[string]interface{}, error)
	CreateSpreadsheetFromTemplateAndData(data []map[string]interface{}, templateFileName string, headerRow int) (*multipart.FileHeader, []byte, error)
}
