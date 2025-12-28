package exceptions

import (
	"testing"
)

func TestOrderNotFoundException_Error(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{"with custom message", "Order with ID 123 not found", "Order with ID 123 not found"},
		{"with empty message", "", "Order not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &OrderNotFoundException{Message: tt.message}
			if err.Error() != tt.expected {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.expected)
			}
		})
	}
}

func TestInvalidOrderDataException_Error(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{"with custom message", "Invalid order ID format", "Invalid order ID format"},
		{"with empty message", "", "Invalid order data"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &InvalidOrderDataException{Message: tt.message}
			if err.Error() != tt.expected {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.expected)
			}
		})
	}
}

func TestAmountNotValidException_Error(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{"with custom message", "Amount must be positive", "Amount must be positive"},
		{"with empty message", "", "Amount is not valid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &AmountNotValidException{Message: tt.message}
			if err.Error() != tt.expected {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.expected)
			}
		})
	}
}

func TestOrderStatusNotFoundException_Error(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{"with custom message", "Status with ID xyz not found", "Status with ID xyz not found"},
		{"with empty message", "", "Order Status not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &OrderStatusNotFoundException{Message: tt.message}
			if err.Error() != tt.expected {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.expected)
			}
		})
	}
}
