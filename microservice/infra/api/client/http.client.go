package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPClient struct {
	baseURL string
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{baseURL: baseURL}
}

func (c *HTTPClient) Get(path string, obj any) error {
	httpClient := new(http.Client)

	url := c.baseURL + path

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("unexpected content type: %s", resp.Header.Get("Content-Type"))
	}

	if err := json.Unmarshal(body, &obj); err != nil {
		return err
	}

	return nil
}
