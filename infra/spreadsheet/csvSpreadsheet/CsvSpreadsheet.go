package csvSpreadsheet

import (
	"encoding/csv"
	"fmt"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"mime/multipart"
	"os"

	
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/envMode"
	"github.com/gtkmk/finder_api/infra/file"
)

type CsvSpreadsheet struct {
	name        string
	fileTempDir string
}

func NewCsvSpreadsheet(name string) port.SpreadsheetInterface {
	return &CsvSpreadsheet{
		name:        name,
		fileTempDir: os.Getenv(envMode.TempDirConst),
	}
}

func (csvSpreadsheet *CsvSpreadsheet) GenerateSpreadsheetFromData(data []map[string]interface{}) (*multipart.FileHeader, []byte, error) {
	fileDir := os.Getenv(envMode.TempDirConst)
	filePath := fmt.Sprintf("%s/%s", fileDir, csvSpreadsheet.name)

	if err := csvSpreadsheet.createCSVFile(data, filePath); err != nil {
		return nil, nil, err
	}

	multipartFile, err := file.GenerateMultipartFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	if err := os.Remove(filePath); err != nil {
		return nil, nil, err
	}

	file, err := file.FileHeaderToBytes(multipartFile)
	if err != nil {
		return nil, nil, err
	}

	return multipartFile, file, nil
}

func (csvSpreadsheet *CsvSpreadsheet) createCSVFile(data []map[string]interface{}, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if len(data) == 0 {
		return helper.ErrorBuilder(helper.InsufficientDateToCreateSpreadsheetConst)
	}

	header := make([]string, 0, len(data[0]))
	for key := range data[0] {
		header = append(header, key)
	}

	err = writer.Write(header)
	if err != nil {
		return err
	}

	for _, row := range data {
		var rowData []string
		for _, key := range header {
			value := fmt.Sprintf("%v", row[key])
			rowData = append(rowData, value)
		}
		err := writer.Write(rowData)
		if err != nil {
			return err
		}
	}

	return nil
}

func (csvSpreadsheet *CsvSpreadsheet) ReadSpreadsheet(filePath string) ([]map[string]interface{}, error) {
	// TODO: implement method ReadSpreadsheet in CsvSpreadsheet when necessary
	return nil, helper.ErrorBuilder(helper.FunctionalityNotImplementedConst)
}

func (csvSpreadsheet *CsvSpreadsheet) CreateSpreadsheetFromTemplateAndData(
	data []map[string]interface{},
	templateFileName string,
	headerRow int,
) (*multipart.FileHeader, []byte, error) {
	return nil, nil, helper.ErrorBuilder(helper.FunctionalityNotImplementedConst)
}
