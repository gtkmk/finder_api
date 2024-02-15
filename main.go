package main

import (
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/database"
	"github.com/gtkmk/finder_api/infra/dotEnv"
	"github.com/gtkmk/finder_api/infra/encryption"
	"github.com/gtkmk/finder_api/infra/envMode"
	"github.com/gtkmk/finder_api/infra/uuid"

	http "github.com/gtkmk/finder_api/adapter/http"
)

func main() {
	if err := defineEnv(); err != nil {
		panic(err)
	}

	connection, idGenerator, passwordEncryptor := defineInfraParams()

	if err := connection.Open(); err != nil {
		panic(err)
	}

	defer connection.Close()

	server := http.NewHttpServer(connection, idGenerator, passwordEncryptor)
	if err := server.Start(); err != nil {
		panic(err)
	}
}

func defineEnv() error {
	dotenv := dotEnv.GodotEnv{}
	if err := dotenv.StartEnv(); err != nil {
		return err
	}

	currentEnvMode := envMode.NewEnvMode()
	if err := currentEnvMode.DefineEnvMode(); err != nil {
		return err
	}

	return nil
}

func defineInfraParams() (
	port.ConnectionInterface,
	port.UuidInterface,
	port.EncryptionInterface,
) {
	connection := database.NewDBConnection()

	idGenerator := uuid.NewUuid()

	passwordEncryptor := encryption.NewPasswordEncryptor()

	return connection, idGenerator, passwordEncryptor
}
