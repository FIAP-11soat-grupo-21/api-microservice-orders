package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPClient(t *testing.T) {
	client := NewHTTPClient("http://localhost:8080")

	assert.NotNil(t, client)
	assert.Equal(t, "http://localhost:8080", client.baseURL)
}

func TestNewHTTPClient_WithDifferentURLs(t *testing.T) {
	testCases := []struct {
		name    string
		baseURL string
	}{
		{"localhost", "http://localhost:8080"},
		{"with https", "https://api.example.com"},
		{"with path", "http://localhost:8080/api/v1"},
		{"empty", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewHTTPClient(tc.baseURL)
			assert.NotNil(t, client)
			assert.Equal(t, tc.baseURL, client.baseURL)
		})
	}
}

func TestHTTPClient_Get_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/products/123", r.URL.Path)

		response := ProductResponseDTO{
			ID:     "123",
			Price:  25.99,
			Active: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	var product ProductResponseDTO
	err := client.Get("/products/123", &product)

	assert.NoError(t, err)
	assert.Equal(t, "123", product.ID)
	assert.Equal(t, 25.99, product.Price)
	assert.True(t, product.Active)
}

func TestHTTPClient_Get_StatusNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Product not found"))
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	var product ProductResponseDTO
	err := client.Get("/products/999", &product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "404")
	assert.Contains(t, err.Error(), "Product not found")
}

func TestHTTPClient_Get_StatusBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	var product ProductResponseDTO
	err := client.Get("/products/invalid", &product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "400")
}

func TestHTTPClient_Get_StatusInternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	var product ProductResponseDTO
	err := client.Get("/products/123", &product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestHTTPClient_Get_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json response"))
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	var product ProductResponseDTO
	err := client.Get("/products/123", &product)

	assert.Error(t, err)
}

func TestHTTPClient_Get_ConnectionError(t *testing.T) {
	client := NewHTTPClient("http://localhost:99999")

	var product ProductResponseDTO
	err := client.Get("/products/123", &product)

	assert.Error(t, err)
}

func TestHTTPClient_Get_InvalidURL(t *testing.T) {
	client := NewHTTPClient("://invalid-url")

	var product ProductResponseDTO
	err := client.Get("/products/123", &product)

	assert.Error(t, err)
}

func TestHTTPClient_Get_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	var product ProductResponseDTO
	err := client.Get("/products/123", &product)

	assert.NoError(t, err)
	assert.Empty(t, product.ID)
	assert.Equal(t, float64(0), product.Price)
	assert.False(t, product.Active)
}

func TestHTTPClient_Get_MultipleRequests(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		response := ProductResponseDTO{
			ID:     "product-" + string(rune('0'+requestCount)),
			Price:  float64(requestCount) * 10.0,
			Active: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	var product1 ProductResponseDTO
	err := client.Get("/products/1", &product1)
	assert.NoError(t, err)

	var product2 ProductResponseDTO
	err = client.Get("/products/2", &product2)
	assert.NoError(t, err)

	assert.Equal(t, 2, requestCount)
}

func TestHTTPClient_Get_LargeResponse(t *testing.T) {
	type LargeResponse struct {
		Items []ProductResponseDTO `json:"items"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		items := make([]ProductResponseDTO, 100)
		for i := 0; i < 100; i++ {
			items[i] = ProductResponseDTO{
				ID:     "item-" + string(rune('0'+i)),
				Price:  float64(i) * 10.0,
				Active: true,
			}
		}
		response := LargeResponse{Items: items}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	var response LargeResponse
	err := client.Get("/products", &response)

	assert.NoError(t, err)
	assert.Len(t, response.Items, 100)
}

func TestHTTPClient_Get_Status300Redirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ProductResponseDTO{
			ID:     "redirect-123",
			Price:  10.00,
			Active: true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMultipleChoices)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	var product ProductResponseDTO
	err := client.Get("/products/123", &product)

	assert.NoError(t, err)
	assert.Equal(t, "redirect-123", product.ID)
}

func TestHTTPClient_Get_Status201Created(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ProductResponseDTO{
			ID:     "new-123",
			Price:  15.99,
			Active: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	var product ProductResponseDTO
	err := client.Get("/products", &product)

	assert.NoError(t, err)
	assert.Equal(t, "new-123", product.ID)
}
