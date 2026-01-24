package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApiClient(t *testing.T) {
	baseURL := "http://localhost:8080"
	client := NewApiClient(baseURL)

	assert.NotNil(t, client)
}

func TestNewApiClient_ReturnsHTTPClient(t *testing.T) {
	baseURL := "http://localhost:8080"
	client := NewApiClient(baseURL)

	// Verify the returned client is an HTTPClient
	httpClient, ok := client.(*HTTPClient)
	assert.True(t, ok, "NewApiClient should return an *HTTPClient")
	assert.NotNil(t, httpClient)
	assert.Equal(t, baseURL, httpClient.baseURL)
}

func TestNewApiClient_ImplementsInterface(t *testing.T) {
	baseURL := "http://localhost:8080"
	client := NewApiClient(baseURL)

	// Verify it implements IApiClient
	var _ IApiClient = client
}

func TestNewApiClient_EmptyBaseURL(t *testing.T) {
	client := NewApiClient("")

	assert.NotNil(t, client)

	httpClient, ok := client.(*HTTPClient)
	assert.True(t, ok)
	assert.Equal(t, "", httpClient.baseURL)
}

func TestNewApiClient_DifferentURLs(t *testing.T) {
	testCases := []struct {
		name    string
		baseURL string
	}{
		{"localhost", "http://localhost:8080"},
		{"production", "https://api.example.com"},
		{"with path", "http://localhost:8080/api/v1"},
		{"ip address", "http://192.168.1.1:3000"},
		{"with port", "http://api.example.com:443"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewApiClient(tc.baseURL)

			assert.NotNil(t, client)

			httpClient, ok := client.(*HTTPClient)
			assert.True(t, ok)
			assert.Equal(t, tc.baseURL, httpClient.baseURL)
		})
	}
}

func TestNewApiClient_MultipleInstances(t *testing.T) {
	client1 := NewApiClient("http://localhost:8080")
	client2 := NewApiClient("http://localhost:9090")

	assert.NotNil(t, client1)
	assert.NotNil(t, client2)

	// They should be different instances
	assert.NotSame(t, client1, client2)

	// With different base URLs
	httpClient1 := client1.(*HTTPClient)
	httpClient2 := client2.(*HTTPClient)
	assert.NotEqual(t, httpClient1.baseURL, httpClient2.baseURL)
}
