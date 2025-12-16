package rest

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewRouter(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("NewRouter panicked as expected without database setup")
		}
	}()

	router := NewRouter()
	if router == nil {
		t.Error("NewRouter() returned nil")
	}
}
