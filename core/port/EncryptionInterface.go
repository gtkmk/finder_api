package port

type EncryptionInterface interface {
	GenerateHashPassword(password string) (string, error)
	CheckHashedPassword(password, hash string) bool
	GenerateRandomPassword() (string, error)
}
