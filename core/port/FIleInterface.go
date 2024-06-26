package port

import "mime/multipart"

type FileInterface interface {
	SetData(data []byte)
	SaveFile() error
	SaveFileFromMultipart(file *multipart.FileHeader, dst string) error
	SaveFromBase64() error
	GetFullPath() string
	RemoveTempFile() error
	DownloadFile() ([]byte, error)
	FileToBase64(filePath string) (string, error)
	GetStaticImageFullPath(imageName string, dst string) string
}
