package file

import (
	"os"

	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/envMode"
	"github.com/gtkmk/finder_api/infra/file/aws"
	"github.com/gtkmk/finder_api/infra/file/toDisk"
)

type FileFactory struct{}

func NewFileFactory() port.FileFactoryInterface {
	return &FileFactory{}
}

func (fileFactory *FileFactory) Make(name string) port.FileInterface {
	if os.Getenv(envMode.ENV_MODE_KEY) == envMode.ENV_MODE_PROD {
		return aws.NewAws(name)
	}

	return toDisk.NewToDisk(name)
}
