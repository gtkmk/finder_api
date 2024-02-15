package port

type ClientInterface interface {
	Execute(method string, path string, body map[string]interface{}, isJson bool) (map[string]interface{}, error)
	CheckAuthorization() (map[string]interface{}, error)
	ExecuteRawDocument(method string, path string, body map[string]interface{}, isJson bool) (string, error)
}
