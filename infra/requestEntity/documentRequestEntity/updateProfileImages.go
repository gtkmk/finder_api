package postRequestEntity

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/core/domain/documentDomain"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/domain/postDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/file"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

const (
	PostMediaConst      = "media"
	CategoryFieldConst  = "o tipo da imagem"
	MediaFileFieldConst = "a imagem da publicação"
)

const (
	MaxTextLengthConst = 800
	MaxLocationLength  = 255
)

var ArrayFileInputName = []string{
	PostMediaConst,
}

var FileTypesToDomain = map[string]string{
	PostMediaConst: documentDomain.PostMediaConst,
}

type UpdateProfileImages struct {
	uuid      port.UuidInterface
	fileUtils port.FileUtilsInterface
	Document  *documentDomain.Document
	Type      string `form:"type" json:"type"`
	UserId    string
	Post      *postDomain.Post
}

func NewUpdateDocumentRequest(context *gin.Context, uuid port.UuidInterface, userId string) (*UpdateProfileImages, error) {
	updateProfileImages := &UpdateProfileImages{
		uuid:      uuid,
		fileUtils: file.NewFileUtils(),
		UserId:    userId,
	}

	if err := context.ShouldBind(updateProfileImages); err != nil {
		return nil, err
	}

	return updateProfileImages, nil
}

func (updateProfileImages *UpdateProfileImages) BuildDocumentObjectAndType(context *gin.Context) (*documentDomain.Document, string, error) {
	form, err := context.MultipartForm()

	if err != nil {
		return nil, updateProfileImages.Type, helper.ErrorBuilder(helper.ErrorGettingFormFileConst, err.Error())
	}

	for _, value := range ArrayFileInputName {
		if formFile, ok := form.File[value]; ok {
			if len(formFile) == 0 {
				continue
			}

			if err := updateProfileImages.logicalConverterIntoDocumentDomain(formFile, context, value); err != nil {
				return nil, updateProfileImages.Type, err
			}
		}
	}

	return updateProfileImages.Document, updateProfileImages.Type, nil
}

func (updateProfileImages *UpdateProfileImages) logicalConverterIntoDocumentDomain(
	file []*multipart.FileHeader,
	context *gin.Context,
	value string,
) error {
	document, convertErr := updateProfileImages.convertMediaIntoDocumentDomain(
		file[0],
		context,
		updateProfileImages.UserId,
		value,
	)

	if convertErr != nil {
		return convertErr
	}

	updateProfileImages.Document = document

	return nil
}

func (updateProfileImages *UpdateProfileImages) convertMediaIntoDocumentDomain(
	file *multipart.FileHeader,
	context *gin.Context,
	documentOwnerId string,
	documentType string,
) (*documentDomain.Document, error) {
	fileExtension := filepath.Ext(file.Filename)

	keySplit := strings.Split(documentType, "_")
	fileType := strings.Join(keySplit[1:], "_")

	fileName := fmt.Sprintf(
		"%s_%s%s",
		updateProfileImages.uuid.GenerateUuid(),
		fileType,
		fileExtension,
	)

	contextFile, err := context.FormFile(documentType)

	if err != nil {
		return nil, helper.ErrorBuilder(helper.IncorrectPasswordOrLoginConst)
	}

	domainDocumentType := translateMimeTypeToDomain(documentType)

	mimeType, err := updateProfileImages.fileUtils.ExtractMimeType(contextFile, fileExtension, domainDocumentType)

	if err != nil {
		return nil, err
	}

	documentId := updateProfileImages.uuid.GenerateUuid()
	document := documentDomain.NewDocument(
		documentId,
		updateProfileImages.Type,
		contextFile,
		fileName,
		nil,
		documentOwnerId,
		mimeType,
		"",
	)

	return document, nil
}

func (updateProfileImages *UpdateProfileImages) Validate(context *gin.Context) error {
	if err := updateProfileImages.validateDocumentType(); err != nil {
		return err
	}

	if err := updateProfileImages.validateMediaField(context); err != nil {
		return err
	}

	return nil
}

func (updateProfileImages *UpdateProfileImages) validateDocumentType() error {
	return requestEntityFieldsValidation.ValidateFieldInArray(
		updateProfileImages.Type,
		CategoryFieldConst,
		documentDomain.AcceptedUserProfileDocumentTypes,
	)
}

func (updateProfileImages *UpdateProfileImages) validateMediaField(context *gin.Context) error {
	form, err := context.MultipartForm()
	if err != nil {
		return helper.ErrorBuilder(helper.ErrorGettingFormFileConst, err.Error())
	}

	if formFile, ok := form.File[PostMediaConst]; ok {
		if len(formFile) == 0 {
			return helper.ErrorBuilder(helper.FieldCannotBeEmptyConst, MediaFileFieldConst)
		}

		for _, file := range formFile {
			if !isImageFile(file) {
				return helper.ErrorBuilder(helper.InvalidFileTypeConst, file.Filename)
			}
		}
	} else {
		return helper.ErrorBuilder(helper.FieldCannotBeEmptyConst, MediaFileFieldConst)
	}

	return nil
}

func isImageFile(file *multipart.FileHeader) bool {
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	ext := filepath.Ext(file.Filename)
	for _, allowedExt := range allowedExtensions {
		if strings.EqualFold(ext, allowedExt) {
			return true
		}
	}
	return false
}

func translateMimeTypeToDomain(documentType string) string {
	documentTypeDomain, found := FileTypesToDomain[documentType]

	if !found {
		return ""
	}

	return documentTypeDomain
}
