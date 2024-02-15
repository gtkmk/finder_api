package dotEnv

import "github.com/joho/godotenv"

type GodotEnv struct{}

func (dotEnv *GodotEnv) StartEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}
