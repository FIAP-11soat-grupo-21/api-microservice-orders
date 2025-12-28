package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewHealthHandler(t *testing.T) {
	handler := NewHealthHandler()

	if handler == nil {
		t.Error("NewHealthHandler() returned nil")
	}
}

func TestHealthHandler_Health(t *testing.T) {
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	handler := NewHealthHandler()
	router.GET("/health", handler.Health)

	req := httptest.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Health() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Health() failed to parse response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Health() status = %v, want healthy", response["status"])
	}
	if response["service"] != "orders-microservice" {
		t.Errorf("Health() service = %v, want orders-microservice", response["service"])
	}
}
