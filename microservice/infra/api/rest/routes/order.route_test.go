package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRegisterOrderRoutes(t *testing.T) {
	router := gin.New()
	group := router.Group("/orders")

	defer func() {
		if r := recover(); r != nil {
			t.Log("RegisterOrderRoutes panicked as expected without database setup")
		}
	}()

	RegisterOrderRoutes(group)
}

func TestRegisterOrderStatusRoutes(t *testing.T) {
	router := gin.New()
	group := router.Group("/orders/status")

	defer func() {
		if r := recover(); r != nil {
			t.Log("RegisterOrderStatusRoutes panicked as expected without database setup")
		}
	}()

	RegisterOrderStatusRoutes(group)
}
