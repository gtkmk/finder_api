package EmptyBodyVerification

import (
	"bytes"
	"github.com/gtkmk/finder_api/core/domain/helper"
	
	"io"
	"net/http"
	"strings"
)

func ValidateBody(req *http.Request) error {
	if req.Body == nil {
		return helper.ErrorBuilder(helper.JsonNotFoundMessageConst)
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		return err
	}

	defer req.Body.Close()

	trimmedBody := strings.TrimSpace(string(body))

	if trimmedBody == "{}" || len(trimmedBody) == 0 {
		return helper.ErrorBuilder(helper.EmptyJsonDataMessageConst)
	}

	req.Body = io.NopCloser(bytes.NewBuffer(body))

	return nil
}
