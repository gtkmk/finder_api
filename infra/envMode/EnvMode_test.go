package envMode

import (
	"os"
	"testing"
)

func TestDefineEnvMode(t *testing.T) {
	envModeSetter := NewEnvMode()

	os.Setenv(ENV_IS_PROD, "true")
	err := envModeSetter.DefineEnvMode()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	os.Setenv(ENV_IS_PROD, "false")
	err = envModeSetter.DefineEnvMode()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	os.Setenv(ENV_IS_PROD, "invalid")
	err = envModeSetter.DefineEnvMode()
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}
