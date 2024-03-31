package file

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gtkmk/finder_api/core/domain/helper"

	"github.com/gtkmk/finder_api/core/domain/documentDomain"
	"github.com/gtkmk/finder_api/core/port"
)

const (
	FileExtensionJpgConst   = ".jpg"
	FileExtensionJpegConst  = ".jpeg"
	FileExtensionPngConst   = ".png"
	FileExtensionGifConst   = ".gif"
	FileExtensionPdfConst   = ".pdf"
	FileExtensionExcelConst = ".xlsx"
)

const (
	MimeTypeJpgConst      = "image/jpeg"
	MimeTypeJpegConst     = "image/jpeg"
	MimeTypePngConst      = "image/png"
	MimeTypeGifConst      = "image/gif"
	MimeTypePdfConst      = "application/pdf"
	MimeTypeExcelConst    = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	MimeTypeNotFoundConst = "application/octet-stream"
)

const (
	MultipartFormMaxMemoryConst = 1024
)

var FileExtensionMimeTypes = map[string]string{
	FileExtensionJpgConst:   MimeTypeJpgConst,
	FileExtensionJpegConst:  MimeTypeJpegConst,
	FileExtensionPngConst:   MimeTypePngConst,
	FileExtensionGifConst:   MimeTypeGifConst,
	FileExtensionPdfConst:   MimeTypePdfConst,
	FileExtensionExcelConst: MimeTypeExcelConst,
}

type FileUtils struct{}

func NewFileUtils() port.FileUtilsInterface {
	return &FileUtils{}
}

func (fileUtils *FileUtils) ExtractMimeType(
	documentContext *multipart.FileHeader,
	fileExtension string,
	documentType string,
) (string, error) {
	fileHeader, err := documentContext.Open()

	if err != nil {
		return "", err
	}

	buffer := make([]byte, 512)

	_, err = fileHeader.Read(buffer)

	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buffer)

	if mimeType == MimeTypeNotFoundConst {
		mimeType, err = fileUtils.ExtractMimeTypeByExtension(
			fileExtension,
			documentType,
		)

		if err != nil {
			return "", err
		}
	}

	if err := fileHeader.Close(); err != nil {
		return "", err
	}

	return mimeType, nil
}

func (fileUtils *FileUtils) ExtractMimeTypeByExtension(
	fileExtension string,
	documentType string,
) (string, error) {
	mimeType, found := FileExtensionMimeTypes[strings.ToLower(fileExtension)]

	if !found {
		documentTypeTranslated, found := documentDomain.DocumentTypesTranslations[documentType]
		if !found {
			return "", fmt.Errorf(helper.ErrorInvalidFileTypeConst, "")
		}

		return "", fmt.Errorf(helper.ErrorInvalidFileTypeConst, documentTypeTranslated)
	}

	return mimeType, nil
}

func (fileUtils *FileUtils) IsValidImageOrPDF(fileHeader *multipart.FileHeader) bool {
	contentType := fileHeader.Header.Get("Content-Type")
	isValidImage := strings.HasPrefix(contentType, "image/") && fileUtils.IsImageExtension(fileHeader.Filename)
	isValidPdf := strings.HasPrefix(contentType, "application/pdf") && fileUtils.IsPDFExtension(fileHeader.Filename)
	return (isValidImage) || (isValidPdf)
}

func (fileUtils *FileUtils) IsImageExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == FileExtensionJpgConst || ext == FileExtensionJpegConst || ext == FileExtensionPngConst || ext == FileExtensionGifConst
}

func (fileUtils *FileUtils) IsPDFExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == FileExtensionPdfConst
}

func (fileUtils *FileUtils) IsValidXLSX(fileHeader *multipart.FileHeader) bool {
	contentType := fileHeader.Header.Get("Content-Type")
	isValidXLSX := strings.HasPrefix(contentType, MimeTypeExcelConst) && fileUtils.IsXLSXExtension(fileHeader.Filename)
	return isValidXLSX
}

func (fileUtils *FileUtils) IsXLSXExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == FileExtensionExcelConst
}

func FileHeaderToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func GenerateMultipartFile(filePath string) (*multipart.FileHeader, error) {
	newFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer newFile.Close()

	var buff bytes.Buffer
	buffWriter := io.Writer(&buff)

	formWriter := multipart.NewWriter(buffWriter)
	formPart, err := formWriter.CreateFormFile("file", filePath)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(formPart, newFile); err != nil {
		return nil, err
	}

	formWriter.Close()

	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())

	multipartForm, err := formReader.ReadForm(MultipartFormMaxMemoryConst)
	if err != nil {
		return nil, err
	}

	files, exists := multipartForm.File["file"]
	if !exists || len(files) == 0 {
		return nil, err
	}

	return files[0], nil
}

func CopyFileIntoAnother(oldFilePath string, newFilePath string) error {
	oldFile, err := os.Open(oldFilePath)
	if err != nil {
		return err
	}

	newFile, err := os.Create(newFilePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(newFile, oldFile)
	if err != nil {
		return err
	}

	if err := oldFile.Close(); err != nil {
		return err
	}

	defer newFile.Close()

	return err
}

func CreateFileHeaderAndBytesByPath(
	filePath string,
) (*multipart.FileHeader, []byte, error) {
	fileHeader, err := GenerateMultipartFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	file, err := FileHeaderToBytes(fileHeader)
	if err != nil {
		return nil, nil, err
	}

	if err := os.Remove(filePath); err != nil {
		return nil, nil, err
	}

	return fileHeader, file, nil
}
