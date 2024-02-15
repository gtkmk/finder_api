package uuid

import uuidGenerator "github.com/google/uuid"

type Uuid struct{}

func NewUuid() *Uuid {
	return &Uuid{}
}

func (uuid *Uuid) GenerateUuid() string {
	Id := uuidGenerator.New().String()
	return Id
}
