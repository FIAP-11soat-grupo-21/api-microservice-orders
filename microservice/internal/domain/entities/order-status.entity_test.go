package entities

import (
	"testing"
)

func TestNewOrderStatus_ValidStatus(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		statusName   string
		expectedName string
	}{
		{"pending status", "status-1", "Pending", "Pending"},
		{"confirmed status", "status-2", "Confirmed", "Confirmed"},
		{"completed status", "status-3", "Completed", "Completed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := NewOrderStatus(tt.id, tt.statusName)
			if err != nil {
				t.Errorf("NewOrderStatus() unexpected error: %v", err)
			}
			if status == nil {
				t.Fatal("NewOrderStatus() returned nil")
			}
			if status.ID != tt.id {
				t.Errorf("NewOrderStatus() ID = %v, want %v", status.ID, tt.id)
			}
			if status.Name.Value() != tt.expectedName {
				t.Errorf("NewOrderStatus() Name = %v, want %v", status.Name.Value(), tt.expectedName)
			}
		})
	}
}

func TestNewOrderStatus_InvalidName(t *testing.T) {
	tests := []struct {
		name       string
		statusName string
	}{
		{"empty name", ""},
		{"too short name", "ab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewOrderStatus("status-1", tt.statusName)
			if err == nil {
				t.Errorf("NewOrderStatus() with name '%v' expected error, got nil", tt.statusName)
			}
		})
	}
}
