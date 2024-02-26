package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"io"
	"net/http"
	"net/url"
)

type HttpClient struct {
	*http.Client
	baseUrl    string
	AuthParams map[string]interface{}
}

const (
	ResponseStatusCodeMinWithoutError = 100
	ResponseStatusCodeMaxWithoutError = 399
)

//TODO - Context: initial class abstraction to be used in the future to consume external APIs

func NewHttpClient(
	baseUrl string,
) *HttpClient {
	return &HttpClient{
		&http.Client{},
		baseUrl,
		nil,
	}
}

func (nativeClient *HttpClient) ExecuteRequest(method string, path string, body map[string]interface{}, isJson bool) (map[string]interface{}, error) {
	resp, err := nativeClient.buildRequest(method, path, body, isJson)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	if err := nativeClient.handleStatusCode(resp.StatusCode, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (nativeClient *HttpClient) ExecuteRawDocumentRequest(method string, path string, body map[string]interface{}, isJson bool) (string, error) {
	resp, err := nativeClient.buildRequest(method, path, body, isJson)

	if err := nativeClient.handleStatusCode(resp.StatusCode, nil); err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	buff := new(bytes.Buffer)
	_, err = buff.ReadFrom(resp.Body)

	if err != nil {
		return "", err
	}

	respBytes := buff.String()

	if err := resp.Body.Close(); err != nil {
		return "", err
	}

	return respBytes, nil
}

func (nativeClient *HttpClient) buildRequest(method string, path string, body map[string]interface{}, isJson bool) (*http.Response, error) {
	baseUrl := fmt.Sprintf(
		"%s/%s",
		nativeClient.baseUrl, path)

	var bodyParams io.Reader

	if body != nil {
		bodyParsed, err := json.Marshal(body)

		if err != nil {
			return nil, err
		}

		bodyParams = bytes.NewBuffer(bodyParsed)
	} else {
		bodyParams = nil
	}

	req, err := http.NewRequest(method, baseUrl, bodyParams)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf(
		"%s",
		"Bearer: 8888888888"),
	)

	if isJson {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := nativeClient.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (nativeClient *HttpClient) handleStatusCode(status int, responseError map[string]interface{}) error {
	if status >= ResponseStatusCodeMinWithoutError && status <= ResponseStatusCodeMaxWithoutError {
		return nil
	}

	return fmt.Errorf(responseError["message"].(string))
}

func (nativeClient *HttpClient) RequestWithJson(method string, path string, data []byte) (map[string]interface{}, error) {
	// TODO
	var response map[string]interface{}

	return response, nil
}

func (nativeClient *HttpClient) RequestWithForm(path string, data url.Values) (map[string]interface{}, error) {
	// TODO
	var res map[string]interface{}

	return res, nil
}

func (nativeClient *HttpClient) Execute(method string, path string, body map[string]interface{}, isJson bool) (map[string]interface{}, error) {
	authParams, _ := nativeClient.CheckAuthorization()
	nativeClient.AuthParams = authParams

	return nativeClient.ExecuteRequest(method, path, body, isJson)
}

func (nativeClient *HttpClient) ExecuteRawDocument(method string, path string, body map[string]interface{}, isJson bool) (string, error) {
	authParams, _ := nativeClient.CheckAuthorization()
	nativeClient.AuthParams = authParams

	return nativeClient.ExecuteRawDocumentRequest(method, path, body, isJson)
}

func (nativeClient *HttpClient) CheckAuthorization() (map[string]interface{}, error) {
	// TODO
	var credentialData []byte

	return nativeClient.RequestWithJson(helper.POST, "urlToService", credentialData)
}
