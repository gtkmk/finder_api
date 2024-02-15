package port

import "mime/multipart"

type FileUtilsInterface interface {
	ExtractMimeType(documentContext *multipart.FileHeader, fileExtension string, documentType string) (string, error)
	IsValidImageOrPDF(fileHeader *multipart.FileHeader) bool
	IsImageExtension(filename string) bool
	IsPDFExtension(filename string) bool
	IsValidXLSX(fileHeader *multipart.FileHeader) bool
	IsXLSXExtension(filename string) bool
}
