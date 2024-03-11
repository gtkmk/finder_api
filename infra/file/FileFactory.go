package file

import (
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/file/toDisk"
)

type FileFactory struct{}

func NewFileFactory() port.FileFactoryInterface {
	return &FileFactory{}
}

func (fileFactory *FileFactory) Make(name string) port.FileInterface {
	return toDisk.NewToDisk(name)
	//return aws.NewAws(name)
}
