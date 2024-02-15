package envMode

import (
	"os"
	"strconv"
)

const (
	ENV_MODE_DEV  = "dev"
	ENV_MODE_PROD = "prod"
	ENV_MODE_KEY  = "ENV_MODE"
	ENV_IS_PROD   = "IS_PROD"
)

type EnvMode struct{}

func NewEnvMode() *EnvMode {
	return &EnvMode{}
}

func (envMode *EnvMode) DefineEnvMode() error {
	isProdEnv := os.Getenv(ENV_IS_PROD)
	isProd, err := strconv.ParseBool(isProdEnv)

	if err != nil {
		return err
	}

	if isProd {
		return envMode.setAsProd()
	} else {
		return envMode.setAsDev()
	}
}

func (envMode *EnvMode) setAsDev() error {
	return os.Setenv(ENV_MODE_KEY, ENV_MODE_DEV)
}

func (envMode *EnvMode) setAsProd() error {
	return os.Setenv(ENV_MODE_KEY, ENV_MODE_PROD)
}
