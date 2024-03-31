package documentDomain

import "mime/multipart"

const (
	PostMediaConst          = "media"
	UserProfilePictureConst = "profile_picture"
	UserProfileBannerConst  = "profile_banner_picture"
)

const (
	PostMediaPortugueseConst          = "Mídia da postagem"
	UserProfilePicturePortugueseConst = "Foto de perfil"
	UserProfileBannerPortugueseConst  = "Banner do perfil"
)

var DocumentTypesTranslations = map[string]string{
	PostMediaConst:          PostMediaPortugueseConst,
	UserProfilePictureConst: UserProfilePicturePortugueseConst,
	UserProfileBannerConst:  UserProfileBannerPortugueseConst,
}

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
