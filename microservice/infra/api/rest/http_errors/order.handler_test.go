package http_errors

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

func TestHandleDomainErrors_InvalidOrderDataException(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	err := &exceptions.InvalidOrderDataException{Message: "Invalid order ID"}
	handled := HandleDomainErrors(err, ctx)

	if !handled {
		t.Error("HandleDomainErrors() should return true for InvalidOrderDataException")
	}
	if w.Code != http.StatusBadRequest {
		t.Errorf("HandleDomainErrors() status = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestHandleDomainErrors_OrderNotFoundException(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	err := &exceptions.OrderNotFoundException{Message: "Order not found"}
	handled := HandleDomainErrors(err, ctx)

	if !handled {
		t.Error("HandleDomainErrors() should return true for OrderNotFoundException")
	}
	if w.Code != http.StatusNotFound {
		t.Errorf("HandleDomainErrors() status = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestHandleDomainErrors_InvalidOrderItemData(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	err := &exceptions.InvalidOrderItemData{Message: "Invalid item data"}
	handled := HandleDomainErrors(err, ctx)

	if !handled {
		t.Error("HandleDomainErrors() should return true for InvalidOrderItemData")
	}
	if w.Code != http.StatusBadRequest {
		t.Errorf("HandleDomainErrors() status = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestHandleDomainErrors_AmountNotValidException(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	err := &exceptions.AmountNotValidException{Message: "Amount must be positive"}
	handled := HandleDomainErrors(err, ctx)

	if !handled {
		t.Error("HandleDomainErrors() should return true for AmountNotValidException")
	}
	if w.Code != http.StatusBadRequest {
		t.Errorf("HandleDomainErrors() status = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestHandleDomainErrors_OrderStatusNotFoundException(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	err := &exceptions.OrderStatusNotFoundException{Message: "Status not found"}
	handled := HandleDomainErrors(err, ctx)

	if !handled {
		t.Error("HandleDomainErrors() should return true for OrderStatusNotFoundException")
	}
	if w.Code != http.StatusNotFound {
		t.Errorf("HandleDomainErrors() status = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestHandleDomainErrors_UnknownError(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	err := errors.New("unknown error")
	handled := HandleDomainErrors(err, ctx)

	if handled {
		t.Error("HandleDomainErrors() should return false for unknown errors")
	}
}
