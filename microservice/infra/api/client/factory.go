package client

func NewApiClient(baseURL string) IApiClient {
	return NewHTTPClient(baseURL)
}
