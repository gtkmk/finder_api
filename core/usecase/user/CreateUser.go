package userUsecase

import (
	"encoding/json"

	"github.com/gtkmk/finder_api/core/domain/documentDomain"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port/sharedMethods"

	"github.com/gtkmk/finder_api/core/domain/notificationDomain"
	"github.com/gtkmk/finder_api/core/domain/userDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type SignUp struct {
	Uuid                port.UuidInterface
	UserDatabase        repositories.UserRepository
	DocumentDatabase    repositories.DocumentRepository
	FileService         port.FileFactoryInterface
	Dist                string
	Transaction         port.ConnectionInterface
	User                *userDomain.User
	NotificationService port.NotificationInterface
	Path                string
	CustomError         port.CustomErrorInterface
	UserEvent           sharedMethods.CreateUserEventInterface
}

func NewCreateUser(
	transaction port.ConnectionInterface,
	uuid port.UuidInterface,
	userDatabase repositories.UserRepository,
	documentDatabase repositories.DocumentRepository,
	fileService port.FileFactoryInterface,
	dist string,
	definedUser *userDomain.User,
	notificationService port.NotificationInterface,
	path string,
	customErrorInterface port.CustomErrorInterface,
	userEvent sharedMethods.CreateUserEventInterface,
) *SignUp {
	return &SignUp{
		Uuid:                uuid,
		Transaction:         transaction,
		UserDatabase:        userDatabase,
		DocumentDatabase:    documentDatabase,
		FileService:         fileService,
		Dist:                dist,
		User:                definedUser,
		NotificationService: notificationService,
		Path:                path,
		CustomError:         customErrorInterface,
		UserEvent:           userEvent,
	}
}

func (signUp *SignUp) Execute(userIP string, userDevice string) error {
	if err := signUp.verifyIfUserExistsByCpf(); err != nil {
		if transactionErr := signUp.Transaction.Rollback(); transactionErr != nil {
			return signUp.CustomError.ThrowError(transactionErr.Error())
		}

		return err
	}

	if err := signUp.verifyIfUserExistsByUserName(); err != nil {
		if transactionErr := signUp.Transaction.Rollback(); transactionErr != nil {
			return signUp.CustomError.ThrowError(transactionErr.Error())
		}

		return err
	}

	if err := signUp.verifyIfUserExists(); err != nil {
		if transactionErr := signUp.Transaction.Rollback(); transactionErr != nil {
			return signUp.CustomError.ThrowError(transactionErr.Error())
		}

		return err
	}

	if err := signUp.createUser(userIP, userDevice); err != nil {
		if transactionErr := signUp.Transaction.Rollback(); transactionErr != nil {
			return signUp.CustomError.ThrowError(transactionErr.Error())
		}

		return err
	}

	if err := signUp.saveDefaultProfileImage(); err != nil {
		if transactionErr := signUp.Transaction.Rollback(); transactionErr != nil {
			return signUp.CustomError.ThrowError(transactionErr.Error())
		}

		return err
	}

	if err := signUp.Transaction.Commit(); err != nil {
		return signUp.CustomError.ThrowError(err.Error())
	}

	_ = signUp.sendSingUpEmail(signUp.User)

	return nil
}

func (signUp *SignUp) verifyIfUserExistsByCpf() error {
	exists := signUp.UserDatabase.VerifyIfUserExistsByCpf(signUp.User.Cpf)

	if exists {
		return signUp.CustomError.ThrowError(helper.UserAlreadyRegisteredConst)
	}

	return nil
}

func (signUp *SignUp) verifyIfUserExistsByUserName() error {
	exists := signUp.UserDatabase.VerifyIfUserExistsByUserName(signUp.User.UserName)

	if exists {
		return signUp.CustomError.ThrowError(helper.UserAlreadyRegisteredConst)
	}

	return nil
}

func (signUp *SignUp) verifyIfUserExists() error {
	dbUser, err := signUp.UserDatabase.FindUserByEmail(signUp.User.Email)

	if err != nil {
		return signUp.CustomError.ThrowError(err.Error())
	}

	if dbUser != nil {
		return signUp.CustomError.ThrowError(helper.UserAlreadyRegisteredWithEmailConst)
	}

	return nil
}

func (signUp *SignUp) createUser(userIP string, userDevice string) error {
	if err := signUp.UserDatabase.CreateUser(signUp.User); err != nil {
		return signUp.CustomError.ThrowError(err.Error())
	}

	if err := signUp.saveNewUserEvent(userIP, userDevice); err != nil {
		return err
	}

	return nil
}

func (signUp *SignUp) saveDefaultProfileImage() error {
	fileService := signUp.FileService.Make("default_user_image")

	defaultProfilePicture := signUp.buildDefaultProfileImageDocumentObj(
		documentDomain.UserProfilePictureConst,
		userDomain.DefaultUserProfileImageConst,
	)

	defaultBannerPicture := signUp.buildDefaultProfileImageDocumentObj(
		documentDomain.UserProfileBannerConst,
		userDomain.DefaultUserBannerImageConst,
	)

	if err := signUp.persistUserDocuments(
		defaultProfilePicture,
		fileService.GetStaticImageFullPath(userDomain.DefaultUserProfileImageConst, signUp.Dist),
	); err != nil {
		return err
	}

	if err := signUp.persistUserDocuments(
		defaultBannerPicture,
		fileService.GetStaticImageFullPath(userDomain.DefaultUserBannerImageConst, signUp.Dist),
	); err != nil {
		return err
	}

	return nil
}

func (signUp *SignUp) buildDefaultProfileImageDocumentObj(documentType string, imageName string) *documentDomain.Document {
	return documentDomain.NewDocument(
		signUp.Uuid.GenerateUuid(),
		documentType,
		nil,
		imageName,
		nil,
		signUp.User.Id,
		documentDomain.DefaultImageMimeTypeConst,
		"",
	)
}

func (signUp *SignUp) persistUserDocuments(document *documentDomain.Document, documentPath string) error {
	if err := signUp.DocumentDatabase.CreateMedia(
		document,
		documentPath,
	); err != nil {
		return signUp.CustomError.ThrowError(err.Error())
	}

	return nil
}

func (signUp *SignUp) sendSingUpEmail(
	User *userDomain.User,
) error {
	notification := notificationDomain.NewNotification(
		userDomain.FirstAccessConst,
		User.Name,
		User.Email,
		"",
		"",
		signUp.Path,
		nil,
	)

	_ = signUp.NotificationService.SendNotifications(notification)

	return nil
}

func (signUp *SignUp) saveNewUserEvent(userIP string, userDevice string) error {
	signUp.User.Password = ""
	interfaceToJson, err := json.MarshalIndent(signUp.User, "", "\t")

	if err != nil {
		return err
	}

	if err := signUp.UserEvent.SaveNewUserEvent(signUp.User.Id, userIP, userDevice, string(interfaceToJson)); err != nil {
		return err
	}
	return nil
}
