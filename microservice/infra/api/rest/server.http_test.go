package rest

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestInit_RouterCreation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := NewRouter()
	
	if router == nil {
		t.Error("Expected router to be created")
	}
	
	routes := router.Routes()
	
	if len(routes) == 0 {
		t.Error("Expected router to have routes configured")
	}
	
	hasHealthRoute := false
	for _, route := range routes {
		if route.Path == "/health" && route.Method == "GET" {
			hasHealthRoute = true
			break
		}
	}
	
	if !hasHealthRoute {
		t.Error("Expected health route to be configured")
	}
}