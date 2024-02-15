package toDisk

import (
	"encoding/base64"
	"fmt"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/envMode"
	"io"
	"mime/multipart"
	"path/filepath"

	"os"
)

type ToDisk struct {
	name        string
	tempFile    *os.File
	data        []byte
	fullPath    string
	tempPath    string
	fileTempDir string
}

func NewToDisk(name string) port.FileInterface {
	return &ToDisk{
		name:        name,
		fileTempDir: os.Getenv(envMode.TempDirConst),
	}
}

func (toDisk *ToDisk) SaveFileFromMultipart(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(fmt.Sprintf("%s/%s", toDisk.fileTempDir, toDisk.name))

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, src)

	toDisk.tempFile = out
	toDisk.fullPath = out.Name()
	toDisk.tempPath = out.Name()

	return err
}

func (toDisk *ToDisk) SaveFile() error {
	tempFile, err := os.Create(fmt.Sprintf("%s/%s", toDisk.fileTempDir, toDisk.name))

	if err != nil {
		return fmt.Errorf("error when creating toDisk: %s", err)
	}

	if _, err := tempFile.Write(toDisk.data); err != nil {
		return fmt.Errorf("erro when write toDisk: %s", err)
	}

	if err := tempFile.Sync(); err != nil {
		return fmt.Errorf("error when commit toDisk: %s", err)
	}

	toDisk.tempFile = tempFile

	toDisk.fullPath = tempFile.Name()

	return nil
}

func (toDisk *ToDisk) SetData(data []byte) {
	toDisk.data = data
}

func (toDisk *ToDisk) GetFullPath() string {
	return toDisk.fullPath
}

func (toDisk *ToDisk) ConvertFromBase64() error {

	decoded, err := base64.StdEncoding.DecodeString(string(toDisk.data))

	if err != nil {
		return fmt.Errorf("error when converting from base64: %s", err)
	}

	toDisk.data = decoded

	return nil
}

func (toDisk *ToDisk) SaveFromBase64() error {
	if err := toDisk.ConvertFromBase64(); err != nil {
		return err
	}

	if err := toDisk.SaveFile(); err != nil {
		return err
	}

	return nil
}
func (toDisk *ToDisk) RemoveTempFile() error {
	return nil
}

func (toDisk *ToDisk) DownloadFile() ([]byte, error) {
	return os.ReadFile(toDisk.name)
}
