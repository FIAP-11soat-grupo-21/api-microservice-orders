package use_cases

import (
	"testing"

	"microservice/internal/domain/entities"
	"microservice/internal/domain/exceptions"
)

func TestDeleteOrderUseCase_ValidatesUUID(t *testing.T) {
	invalidIDs := []string{
		"",
		"invalid-uuid",
		"123",
		"not-a-uuid-at-all",
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

func TestDeleteOrderUseCase_AcceptsValidUUID(t *testing.T) {
	validID := "550e8400-e29b-41d4-a716-446655440000"
	err := entities.ValidateID(validID)
	if err != nil {
		t.Errorf("ValidateID(%v) unexpected error: %v", validID, err)
	}
}

func TestDeleteOrderUseCase_OrderNotFoundExceptionMessage(t *testing.T) {
	err := &exceptions.OrderNotFoundException{}
	expectedMsg := "Order not found"
	if err.Error() != expectedMsg {
		t.Errorf("OrderNotFoundException.Error() = %v, want %v", err.Error(), expectedMsg)
	}
}

func TestDeleteOrderUseCase_OrderNotFoundExceptionCustomMessage(t *testing.T) {
	customMsg := "Order with ID 123 not found"
	err := &exceptions.OrderNotFoundException{Message: customMsg}
	if err.Error() != customMsg {
		t.Errorf("OrderNotFoundException.Error() = %v, want %v", err.Error(), customMsg)
	}
}

func TestNewDeleteOrderUseCase(t *testing.T) {
	_ = NewDeleteOrderUseCase
}
func TestDeleteOrderUseCase_NewDeleteOrderUseCase(t *testing.T) {
	mockGateway := NewMockOrderGateway()
	uc := NewDeleteOrderUseCase(mockGateway)

	if uc == nil {
		t.Error("Expected use case to be created")
	}
}

func TestDeleteOrderUseCase_Execute_InvalidID(t *testing.T) {
	mockGateway := NewMockOrderGateway()
	uc := NewDeleteOrderUseCase(mockGateway)

	err := uc.Execute("invalid-id")
	if err == nil {
		t.Error("Expected error for invalid ID")
	}

	if _, ok := err.(*exceptions.InvalidOrderDataException); !ok {
		t.Errorf("Expected InvalidOrderDataException, got %T", err)
	}
}

func TestDeleteOrderUseCase_Execute_EmptyID(t *testing.T) {
	mockGateway := NewMockOrderGateway()
	uc := NewDeleteOrderUseCase(mockGateway)

	err := uc.Execute("")
	if err == nil {
		t.Error("Expected error for empty ID")
	}
}

func TestDeleteOrderUseCase_Execute_ValidIDFormat(t *testing.T) {
	validID := "550e8400-e29b-41d4-a716-446655440000"
	err := entities.ValidateID(validID)
	if err != nil {
		t.Errorf("Expected no error for valid UUID, got %v", err)
	}
}