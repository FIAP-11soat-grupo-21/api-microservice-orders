package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"microservice/internal/domain/exceptions"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestErrorHandlerMiddleware_NoErrors(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, router := gin.CreateTestContext(w)

	router.Use(ErrorHandlerMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	ctx.Request = httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, ctx.Request)

	if w.Code != http.StatusOK {
		t.Errorf("ErrorHandlerMiddleware() status = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestErrorHandlerMiddleware_DomainError(t *testing.T) {
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.Use(ErrorHandlerMiddleware())
	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(&exceptions.OrderNotFoundException{Message: "Order not found"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("ErrorHandlerMiddleware() status = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestErrorHandlerMiddleware_InvalidOrderDataException(t *testing.T) {
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.Use(ErrorHandlerMiddleware())
	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(&exceptions.InvalidOrderDataException{Message: "Invalid data"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("ErrorHandlerMiddleware() status = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestErrorHandlerMiddleware_UnknownError(t *testing.T) {
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.Use(ErrorHandlerMiddleware())
	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(errors.New("unknown error"))
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("ErrorHandlerMiddleware() status = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}
