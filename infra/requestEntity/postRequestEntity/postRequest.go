package postRequestEntity

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/documentDomain"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/domain/postDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/file"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

const (
	PostMediaConst       = "media"
	TextFieldConst       = "o texto da publicação"
	LocationFieldConst   = "a localização da publicação"
	RewardFieldConst     = "a recompensa da publicação"
	PrivacyFieldConst    = "a deinição de privacidade da publicação"
	LostFoundFieldConst  = "o tipo de achado/perdido da publicação"
	AnimalTypeFieldConst = "o tipo de animal da publicação"
	AnimalSizeFieldConst = "o porte do animal da publicação"
	CategoryFieldConst   = "a categoria da publicação"
	MediaFileFieldConst  = "a imagem da publicação"
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

type PostRequest struct {
	uuid       port.UuidInterface
	fileUtils  port.FileUtilsInterface
	Document   *documentDomain.Document
	PostId     string  `form:"post_id" json:"post_id"`
	Text       string  `form:"text" json:"text"`
	Location   string  `form:"location" json:"location"`
	Reward     bool    `form:"reward" json:"reward"`
	LostFound  *string `form:"lost_found" json:"lost_found"`
	Privacy    string  `form:"privacy" json:"privacy"`
	Category   string  `form:"category" json:"category"`
	AnimalType *string `form:"animal_type" json:"animal_type"`
	AnimalSize *string `form:"animal_size" json:"animal_size"`
	UserId     string
	Post       *postDomain.Post
}

func NewPostRequest(context *gin.Context, uuid port.UuidInterface, userId string) (*PostRequest, error) {
	postRequest := &PostRequest{
		uuid:      uuid,
		fileUtils: file.NewFileUtils(),
		UserId:    userId,
	}

	if err := context.ShouldBind(postRequest); err != nil {
		return nil, err
	}

	return postRequest, nil
}

func (postRequest *PostRequest) IterateIntoFiles(context *gin.Context) error {
	form, err := context.MultipartForm()

	if err != nil {
		return helper.ErrorBuilder(helper.ErrorGettingFormFileConst, err.Error())
	}

	for _, value := range ArrayFileInputName {
		if formFile, ok := form.File[value]; ok {
			if len(formFile) == 0 {
				continue
			}

			if err := postRequest.logicalConverterIntoDocumentDomain(formFile, context, value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (postRequest *PostRequest) logicalConverterIntoDocumentDomain(
	file []*multipart.FileHeader,
	context *gin.Context,
	value string,
) error {
	document, convertErr := postRequest.convertPostMediaIntoDocumentDomain(
		file[0],
		context,
		postRequest.UserId,
		value,
	)

	if convertErr != nil {
		return convertErr
	}

	postRequest.Document = document

	return nil
}

func (postRequest *PostRequest) convertPostMediaIntoDocumentDomain(
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
		postRequest.uuid.GenerateUuid(),
		fileType,
		fileExtension,
	)

	contextFile, err := context.FormFile(documentType)

	if err != nil {
		return nil, helper.ErrorBuilder(helper.IncorrectPasswordOrLoginConst)
	}

	domainDocumentType := translateMimeTypeToDomain(documentType)

	mimeType, err := postRequest.fileUtils.ExtractMimeType(contextFile, fileExtension, domainDocumentType)

	if err != nil {
		return nil, err
	}

	documentId := postRequest.uuid.GenerateUuid()
	postRequest.assignNewPostId()
	document := documentDomain.NewDocument(
		documentId,
		fileType,
		contextFile,
		fileName,
		&postRequest.PostId,
		documentOwnerId,
		mimeType,
		"",
	)

	return document, nil
}

func (postRequest *PostRequest) assignNewPostId() {
	if postRequest.PostId == "" {
		postRequest.PostId = postRequest.uuid.GenerateUuid()
	}
}

func (postRequest *PostRequest) Validate(context *gin.Context, edition bool) error {
	if err := postRequest.validatePostFields(); err != nil {
		return err
	}

	if !edition {
		if err := postRequest.validateMediaField(context); err != nil {
			return err
		}
	}

	return nil
}

func (postRequest *PostRequest) validatePostFields() error {
	if err := postRequest.validateText(); err != nil {
		return err
	}

	if err := postRequest.validateLocation(); err != nil {
		return err
	}

	if err := postRequest.validatePrivacy(); err != nil {
		return err
	}

	if err := postRequest.validateLostFound(); err != nil {
		return err
	}

	if err := postRequest.validateAnimalFields(); err != nil {
		return err
	}

	if err := postRequest.validateCategory(); err != nil {
		return err
	}

	return nil
}

func (postRequest *PostRequest) validateText() error {
	return requestEntityFieldsValidation.ValidateField(
		postRequest.Text,
		TextFieldConst,
		MaxTextLengthConst,
	)
}

func (postRequest *PostRequest) validateLocation() error {
	return requestEntityFieldsValidation.ValidateField(
		postRequest.Location,
		LocationFieldConst,
		MaxLocationLength,
	)
}

func (postRequest *PostRequest) validatePrivacy() error {
	return requestEntityFieldsValidation.ValidateFieldInArray(
		postRequest.Privacy,
		PrivacyFieldConst,
		postDomain.AcceptedPrivacySettings,
	)
}

func (postRequest *PostRequest) validateLostFound() error {
	if len(*postRequest.LostFound) == 0 {
		postRequest.LostFound = nil
	}

	if postRequest.LostFound != nil {
		return requestEntityFieldsValidation.ValidateFieldInArray(
			*postRequest.LostFound,
			LostFoundFieldConst,
			postDomain.LostAndFoundStatus,
		)
	}
	return nil
}

func (postRequest *PostRequest) validateAnimalFields() error {
	if postRequest.LostFound != nil {
		if postRequest.AnimalType == nil {
			return helper.ErrorBuilder(
				helper.FieldCannotBeEmptyConst,
				AnimalTypeFieldConst,
			)
		}

		if postRequest.AnimalSize == nil {
			return helper.ErrorBuilder(
				helper.FieldCannotBeEmptyConst,
				AnimalSizeFieldConst,
			)
		}
	}

	if postRequest.AnimalType != nil || *postRequest.AnimalType != "" {
		if err := requestEntityFieldsValidation.ValidateFieldInArray(
			*postRequest.AnimalType,
			AnimalTypeFieldConst,
			postDomain.AcceptedAnimalTypes,
		); err != nil {
			return err
		}
	}

	if postRequest.AnimalSize != nil || *postRequest.AnimalSize != "" {
		if err := requestEntityFieldsValidation.ValidateFieldInArray(
			*postRequest.AnimalSize,
			AnimalSizeFieldConst,
			postDomain.AcceptedAnimalSizes,
		); err != nil {
			return err
		}
	}

	return nil
}

func (postRequest *PostRequest) validateCategory() error {
	return requestEntityFieldsValidation.ValidateFieldInArray(
		postRequest.Category,
		CategoryFieldConst,
		postDomain.AcceptedCategories,
	)
}

func (postRequest *PostRequest) validateMediaField(context *gin.Context) error {
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

func (postRequest *PostRequest) BuildPostObject(postId *string) (*postDomain.Post, error) {
	dateTime, err := datetimeDomain.CreateNow()
	if err != nil {
		return nil, err
	}

	if postId == nil {
		postId = &postRequest.PostId
	}

	return postDomain.NewPost(
		*postId,
		postRequest.Text,
		postRequest.Document,
		postRequest.Location,
		postRequest.Reward,
		postRequest.Privacy,
		0,
		postRequest.Category,
		postRequest.LostFound,
		postRequest.AnimalType,
		postRequest.AnimalSize,
		false,
		nil,
		postRequest.UserId,
		&dateTime,
		nil,
	), nil
}
