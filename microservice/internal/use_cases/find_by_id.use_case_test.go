package use_cases

import (
	"testing"
	"time"

	"microservice/internal/domain/entities"
	"microservice/internal/domain/exceptions"
)

func TestFindByIDUseCase_ValidatesUUID(t *testing.T) {
	invalidIDs := []string{
		"",
		"invalid-uuid",
		"123",
		"not-a-uuid",
		"550e8400-e29b-41d4-a716",
	}

	for _, id := range invalidIDs {
		t.Run(id, func(t *testing.T) {
			err := entities.ValidateID(id)
			if err == nil {
				t.Errorf("ValidateID(%v) expected error, got nil", id)
			}
			if _, ok := err.(*exceptions.InvalidOrderDataException); !ok {
				t.Errorf("ValidateID(%v) expected InvalidOrderDataException, got %T", id, err)
			}
		})
	}
}

func TestFindByIDUseCase_AcceptsValidUUID(t *testing.T) {
	validIDs := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"f47ac10b-58cc-4372-a567-0e02b2c3d479",
	}

	for _, id := range validIDs {
		t.Run(id, func(t *testing.T) {
			err := entities.ValidateID(id)
			if err != nil {
				t.Errorf("ValidateID(%v) unexpected error: %v", id, err)
			}
		})
	}
}

func TestFindByIDUseCase_OrderNotFoundExceptionMessage(t *testing.T) {
	err := &exceptions.OrderNotFoundException{}
	if err.Error() != "Order not found" {
		t.Errorf("OrderNotFoundException.Error() = %v, want 'Order not found'", err.Error())
	}

	errWithMsg := &exceptions.OrderNotFoundException{Message: "Order with ID xyz not found"}
	if errWithMsg.Error() != "Order with ID xyz not found" {
		t.Errorf("OrderNotFoundException.Error() = %v, want 'Order with ID xyz not found'", errWithMsg.Error())
	}
}

func TestNewFindOrderByIDUseCase(t *testing.T) {
	_ = NewFindOrderByIDUseCase
}
func TestFindOrderByIDUseCase_NewFindOrderByIDUseCase(t *testing.T) {
	mockGateway := NewMockOrderGateway()
	uc := NewFindOrderByIDUseCase(mockGateway)

	if uc == nil {
		t.Error("Expected use case to be created")
	}
}

func TestFindOrderByIDUseCase_Execute_InvalidID(t *testing.T) {
	mockGateway := NewMockOrderGateway()
	uc := NewFindOrderByIDUseCase(mockGateway)

	_, err := uc.Execute("invalid-id")
	if err == nil {
		t.Error("Expected error for invalid ID")
	}

	if _, ok := err.(*exceptions.InvalidOrderDataException); !ok {
		t.Errorf("Expected InvalidOrderDataException, got %T", err)
	}
}

func TestFindOrderByIDUseCase_Execute_EmptyID(t *testing.T) {
	mockGateway := NewMockOrderGateway()
	uc := NewFindOrderByIDUseCase(mockGateway)

	_, err := uc.Execute("")
	if err == nil {
		t.Error("Expected error for empty ID")
	}
}

func TestFindOrderByIDUseCase_Execute_ValidIDFormat(t *testing.T) {
	validID := "550e8400-e29b-41d4-a716-446655440000"
	err := entities.ValidateID(validID)
	if err != nil {
		t.Errorf("Expected no error for valid UUID, got %v", err)
	}
}

func TestFindOrderByIDUseCase_Execute_ReturnsEmptyOrderOnError(t *testing.T) {
	mockGateway := NewMockOrderGateway()
	uc := NewFindOrderByIDUseCase(mockGateway)

	order, err := uc.Execute("invalid-id")
	if err == nil {
		t.Error("Expected error for invalid ID")
	}

	if !order.IsEmpty() {
		t.Error("Expected empty order on error")
	}
}

// Comprehensive tests using mocks for full coverage

func TestFindOrderByIDUseCase_Execute_Success(t *testing.T) {
	mockGateway := NewMockOrderGateway()
	uc := NewFindOrderByIDUseCase(mockGateway)

	// Create and add a test order
	validID := "550e8400-e29b-41d4-a716-446655440000"
	customerID := "customer-123"
	status, _ := entities.NewOrderStatus("pending", "Pending")
	expectedOrder, _ := entities.NewOrderWithItems(validID, &customerID, 25.0, *status, []entities.OrderItem{}, time.Now(), nil)
	mockGateway.AddOrder(expectedOrder)

	order, err := uc.Execute(validID)
	if err != nil {
		t.Errorf("Expected no error for successful find, got %v", err)
	}

	if order.ID != validID {
		t.Errorf("Expected order ID %s, got %s", validID, order.ID)
	}

	if order.IsEmpty() {
		t.Error("Expected non-empty order")
	}
}

func TestFindOrderByIDUseCase_Execute_OrderNotFound(t *testing.T) {
	mockGateway := NewMockOrderGateway()
	uc := NewFindOrderByIDUseCase(mockGateway)

	validID := "550e8400-e29b-41d4-a716-446655440000"

	order, err := uc.Execute(validID)
	if err == nil {
		t.Error("Expected error when order not found")
	}

	if _, ok := err.(*exceptions.OrderNotFoundException); !ok {
		t.Errorf("Expected OrderNotFoundException, got %T", err)
	}

	if !order.IsEmpty() {
		t.Error("Expected empty order when not found")
	}
}

func TestFindOrderByIDUseCase_Execute_GatewayError(t *testing.T) {
	mockGateway := NewMockOrderGateway()
	mockGateway.SetShouldFailFindByID(true)
	uc := NewFindOrderByIDUseCase(mockGateway)

	validID := "550e8400-e29b-41d4-a716-446655440000"

	order, err := uc.Execute(validID)
	if err == nil {
		t.Error("Expected error when gateway fails")
	}

	if _, ok := err.(*exceptions.OrderNotFoundException); !ok {
		t.Errorf("Expected OrderNotFoundException, got %T", err)
	}

	if !order.IsEmpty() {
		t.Error("Expected empty order on gateway error")
	}
}

func TestFindOrderByIDUseCase_Execute_IDValidation(t *testing.T) {
	testCases := []struct {
		name  string
		id    string
		valid bool
	}{
		{"valid UUID", "550e8400-e29b-41d4-a716-446655440000", true},
		{"empty ID", "", false},
		{"invalid format", "invalid-id", false},
		{"short UUID", "550e8400-e29b", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := entities.ValidateID(tc.id)
			if tc.valid && err != nil {
				t.Errorf("Expected no error for valid case, got %v", err)
			}
			if !tc.valid && err == nil {
				t.Error("Expected error for invalid case")
			}
		})
	}
}

func TestFindOrderByIDUseCase_Execute_ErrorHandling(t *testing.T) {
	orderNotFoundErr := &exceptions.OrderNotFoundException{}
	if orderNotFoundErr.Error() != "Order not found" {
		t.Errorf("Expected 'Order not found', got '%s'", orderNotFoundErr.Error())
	}

	customErr := &exceptions.OrderNotFoundException{Message: "Custom message"}
	if customErr.Error() != "Custom message" {
		t.Errorf("Expected 'Custom message', got '%s'", customErr.Error())
	}
}
