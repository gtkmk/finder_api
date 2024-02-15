package xlsxSpreadsheet

import (
	"fmt"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"mime/multipart"
	"path/filepath"

	
	"github.com/gtkmk/finder_api/infra/file"
)

type XlsxSpreadsheet struct {
	name string
}

func NewXlsxSpreadsheet(name string) *XlsxSpreadsheet {
	return &XlsxSpreadsheet{
		name: name,
	}
}

func (xlsxSpreadsheet *XlsxSpreadsheet) GenerateSpreadsheetFromData(data []map[string]interface{}) (*multipart.FileHeader, []byte, error) {
	// TODO: implement method GenerateSpreadsheetFromData in XlsxSpreadsheet when necessary
	return nil, nil, helper.ErrorBuilder(helper.FunctionalityNotImplementedConst)
}

func (xlsxSpreadsheet *XlsxSpreadsheet) ReadSpreadsheet(filePath string) ([]map[string]interface{}, error) {
	xlsx, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	data := make([]map[string]interface{}, 0)

	columnIndices := make(map[string]int, len(rows[0]))
	for colIndex, colName := range rows[0] {
		columnIndices[colName] = colIndex
	}

	for _, row := range rows[1:] {
		rowData := make(map[string]interface{}, len(rows[0]))
		for colName, colIndex := range columnIndices {
			if colIndex >= len(row) {
				rowData[colName] = nil
			} else {
				rowData[colName] = row[colIndex]
			}
		}

		data = append(data, rowData)
	}

	return data, nil
}

func (xlsxSpreadsheet *XlsxSpreadsheet) CreateSpreadsheetFromTemplateAndData(
	data []map[string]interface{},
	templateFileName string,
	headerRow int,
) (*multipart.FileHeader, []byte, error) {
	templateFilePath := filepath.Join("infra/templates", templateFileName)
	newFilePath := filepath.Join("infra/templates", xlsxSpreadsheet.name)

	if err := file.CopyFileIntoAnother(
		templateFilePath,
		newFilePath,
	); err != nil {
		return nil, nil, err
	}

	if err := xlsxSpreadsheet.openFileAndInsertRowsByData(
		data,
		newFilePath,
		headerRow,
	); err != nil {
		return nil, nil, err
	}

	fileHeader, file, err := file.CreateFileHeaderAndBytesByPath(newFilePath)
	if err != nil {
		return nil, nil, err
	}

	return fileHeader, file, nil
}

func (xlsxSpreadsheet *XlsxSpreadsheet) openFileAndInsertRowsByData(
	data []map[string]interface{},
	filePath string,
	headerRow int,
) error {
	xlsx, err := excelize.OpenFile(filepath.Join(filePath))
	if err != nil {
		return err
	}

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return err
	}

	columnNames := rows[headerRow]

	for i, rowData := range data {
		row := []string{}
		for _, colName := range columnNames {
			value, found := rowData[colName]
			if found {
				row = append(row, fmt.Sprint(value))
			} else {
				row = append(row, "")
			}
		}

		if err := xlsx.SetSheetRow(sheetName, fmt.Sprintf("A%d", len(rows)+i+1), &row); err != nil {
			return err
		}
	}

	err = xlsx.Save()
	if err != nil {
		return err
	}

	return nil
}
