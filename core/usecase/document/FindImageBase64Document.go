package documentUsecase

import (
	"github.com/gtkmk/finder_api/core/port"
)

type FindDocumentImageBase64 struct {
	FileService port.FileFactoryInterface
	customError port.CustomErrorInterface
}

func NewFindDocumentImageBase64(
	fileService port.FileFactoryInterface,
	customError port.CustomErrorInterface,
) *FindDocumentImageBase64 {
	return &FindDocumentImageBase64{
		FileService: fileService,
		customError: customError,
	}
}

func (findDocumentImageBase64 *FindDocumentImageBase64) Execute(documentPath string) (string, error) {
	fileService := findDocumentImageBase64.FileService.Make(documentPath)

	return fileService.FileToBase64(documentPath)
}
