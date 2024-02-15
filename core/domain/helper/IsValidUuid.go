package helper

import (
	"github.com/google/uuid"
)

func IsValidUUID(hash string) bool {
	_, err := uuid.Parse(hash)
	return err == nil
}
