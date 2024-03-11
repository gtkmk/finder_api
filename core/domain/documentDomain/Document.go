package documentDomain

import "mime/multipart"

const (
	CsvMimeTypeConst         = "text/csv"
	CsvFileExtensionConst    = ".csv"
	SpreadsheetTypeCsvConst  = "csv"
	SpreadsheetTypeXlsxConst = "xlsx"
	XslxMimeTypeConst        = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	XlsxFileExtensionConst   = ".xslx"
)

const (
	PostConst           = "identity_front"
	ProfilePictureConst = "profile_picture"
	ProfileBannerConst  = "profile_banner"
)

type Document struct {
	ID            string
	Type          string
	File          *multipart.FileHeader
	NewName       string
	PostId        string
	OwnerId       string
	MimeType      string
	FileExtension string
	Data          []byte
}

func NewDocument(
	id string,
	documentType string,
	documentFile *multipart.FileHeader,
	newName string,
	postId string,
	ownerId string,
	mimeType string,
	fileExtension string,
) *Document {
	return &Document{
		ID:            id,
		Type:          documentType,
		File:          documentFile,
		NewName:       newName,
		PostId:        postId,
		OwnerId:       ownerId,
		MimeType:      mimeType,
		FileExtension: fileExtension,
	}
}
