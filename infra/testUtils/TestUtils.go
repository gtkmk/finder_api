package testUtils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func EmptyField(fieldName string) string {
	return fmt.Sprintf("%s não pode ser vazio.", fieldName)
}

func InvalidUUID(fieldName string) string {
	return fmt.Sprintf("É necessário informar %s.", fieldName)
}

func MaximumLengthErrorValidationUtil(fieldName string, maxLength int) string {
	return fmt.Sprintf("%s não pode ter mais de %d caracteres.", fieldName, maxLength)
}

func WriteMultipartFormFieldHelper(headerContent map[string]interface{}, writer *multipart.Writer, t *testing.T) {
	for key, value := range headerContent {
		formFieldWriter, err := writer.CreateFormField(key)

		assert.Nil(t, err)

		_, err = formFieldWriter.Write([]byte(value.(string)))

		assert.Nil(t, err)
	}
}

func WriteMultipartFormFileHelper(headerContent map[string]interface{}, fileHeader textproto.MIMEHeader, writer *multipart.Writer, t *testing.T) {
	for key, value := range headerContent {
		fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, key, value))
		fileHeader.Set("Content-Type", "image/")
		writer, err := writer.CreatePart(fileHeader)
		assert.Nil(t, err)

		file, err := os.Open(fmt.Sprintf("multipartTests/%s", value.(string)))
		assert.Nil(t, err)
		_, err = io.Copy(writer, file)

		assert.Nil(t, err)

		err = file.Close()
		assert.Nil(t, err)
	}
}
