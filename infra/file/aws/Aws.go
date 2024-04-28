package aws

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gtkmk/finder_api/infra/envMode"

	awspackage "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gtkmk/finder_api/core/port"
)

type Aws struct {
	name        string
	tempFile    *os.File
	data        []byte
	fullPath    string
	tempPath    string
	fileTempDir string
	bucket      string
	region      string
}

func NewAws(name string) port.FileInterface {
	bucket := os.Getenv(envMode.AwsBucketConst)
	region := os.Getenv(envMode.AwsRegionConst)

	return &Aws{
		name:        name,
		fileTempDir: os.Getenv(envMode.TempDirConst),
		bucket:      bucket,
		region:      region,
	}
}

func (aws *Aws) defineSession() (*session.Session, error) {
	return session.NewSession(&awspackage.Config{
		Region:      awspackage.String(aws.region),
		Credentials: credentials.NewEnvCredentials(),
	})
}

func (aws *Aws) SaveFileFromMultipart(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(fmt.Sprintf("%s/%s", aws.fileTempDir, aws.name))

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, src)

	if err != nil {
		return err
	}

	filetemp, err := os.Open(out.Name())

	if err != nil {
		return err
	}

	aws.tempPath = out.Name()

	return aws.uploadToS3(filetemp)
}

func (aws *Aws) uploadToS3(tempFile *os.File) error {
	sessionConst, err := aws.defineSession()

	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sessionConst)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: awspackage.String(aws.bucket),
		Key:    awspackage.String(aws.name),
		Body:   tempFile,
	})

	if err != nil {
		return err
	}

	aws.tempFile = tempFile
	aws.fullPath = result.Location

	return err
}

func (aws *Aws) SaveFile() error {
	fileTempDir := os.Getenv(envMode.TempDirConst)
	tempFile, err := os.Create(fmt.Sprintf("%s/%s", fileTempDir, aws.name))

	if err != nil {
		return fmt.Errorf("error when creating toDisk: %s", err)
	}

	if _, err := tempFile.Write(aws.data); err != nil {
		return fmt.Errorf("erro when write toDisk: %s", err)
	}

	if err := tempFile.Sync(); err != nil {
		return fmt.Errorf("error when commit toDisk: %s", err)
	}

	fileTemp, err := os.Open(tempFile.Name())

	if err != nil {
		log.Println("Unable to open file", err)
	}

	defer fileTemp.Close()

	return aws.uploadToS3(fileTemp)
}

func (aws *Aws) SetData(data []byte) {
	aws.data = data
}

func (aws *Aws) GetFullPath() string {
	return aws.fullPath
}

func (aws *Aws) ConvertFromBase64() error {
	decoded, err := base64.StdEncoding.DecodeString(string(aws.data))

	if err != nil {
		return fmt.Errorf("error when converting from base64: %s", err)
	}

	aws.data = decoded

	return nil
}

func (aws *Aws) SaveFromBase64() error {
	if err := aws.ConvertFromBase64(); err != nil {
		return err
	}

	if err := aws.SaveFile(); err != nil {
		return err
	}

	return nil
}

func (aws *Aws) RemoveTempFile() error {
	return os.Remove(aws.tempPath)
}

func (aws *Aws) DownloadFile() ([]byte, error) {
	sess, err := aws.defineSession()
	if err != nil {
		return nil, err
	}

	s3Service := s3.New(sess)

	fileName := formatFileName(aws.name)
	filePath := filepath.Join("/tmp", fileName)

	if err := downloadS3File(s3Service, aws.bucket, fileName, filePath); err != nil {
		return nil, err
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := os.Remove(filePath); err != nil {
		return nil, err
	}

	return fileContent, nil
}

func formatFileName(url string) string {
	parts := strings.Split(url, "/")
	fileName := parts[len(parts)-1]
	return path.Base(fileName)
}

func downloadS3File(svc *s3.S3, bucket, fileName, filePath string) error {
	input := &s3.GetObjectInput{
		Bucket: awspackage.String(bucket),
		Key:    awspackage.String(fileName),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, result.Body)
	if err != nil {
		return err
	}

	return nil
}

func (aws *Aws) FileToBase64(filePath string) (string, error) {
	// TODO: Implement this

	return "", nil
}
