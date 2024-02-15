package encryption

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gtkmk/finder_api/core/port"
	"golang.org/x/crypto/bcrypt"
)

type PasswordEncryptor struct{}

const (
	CharsetConst              = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	RandomPasswordLengthConst = 10
	HashSaltConst             = 13
)

func NewPasswordEncryptor() port.EncryptionInterface {
	return &PasswordEncryptor{}
}

func (passwordEncryptor *PasswordEncryptor) GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), HashSaltConst)
	return string(bytes), err
}

func (passwordEncryptor *PasswordEncryptor) CheckHashedPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (passwordEncryptor *PasswordEncryptor) GenerateRandomPassword() (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwordEncryptor.generateRandomString()), HashSaltConst)

	if err != nil {
		return "", fmt.Errorf("error in password hash")
	}

	return string(bytes), err
}

func (passwordEncryptor *PasswordEncryptor) generateRandomString() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, RandomPasswordLengthConst)
	for i := range result {
		result[i] = CharsetConst[seededRand.Intn(len(CharsetConst))]
	}

	return string(result)
}
