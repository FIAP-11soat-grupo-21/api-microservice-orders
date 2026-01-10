package use_cases

import (
	"testing"

	"microservice/internal/adapters/dtos"
	"microservice/internal/domain/entities"
)

func TestFindAllStatusUseCase_OrderStatusCreation(t *testing.T) {
	status, err := entities.NewOrderStatus("status-1", "Pending")
	if err != nil {
		t.Errorf("NewOrderStatus() unexpected error: %v", err)
	}
	if status.ID != "status-1" {
		t.Errorf("OrderStatus.ID = %v, want status-1", status.ID)
	}
	if status.Name.Value() != "Pending" {
		t.Errorf("OrderStatus.Name = %v, want Pending", status.Name.Value())
	}
}

func TestFindAllStatusUseCase_OrderStatusDTOStructure(t *testing.T) {
	dto := dtos.OrderStatusDTO{
		ID:   "status-1",
		Name: "Pending",
	}

	if dto.ID != "status-1" {
		t.Errorf("OrderStatusDTO.ID = %v, want status-1", dto.ID)
	}
	if dto.Name != "Pending" {
		t.Errorf("OrderStatusDTO.Name = %v, want Pending", dto.Name)
	}
}

func TestFindAllStatusUseCase_OrderStatusResponseDTOStructure(t *testing.T) {
	dto := dtos.OrderStatusResponseDTO{
		ID:   "status-1",
		Name: "Pending",
	}

	if dto.ID != "status-1" {
		t.Errorf("OrderStatusResponseDTO.ID = %v, want status-1", dto.ID)
	}
	if dto.Name != "Pending" {
		t.Errorf("OrderStatusResponseDTO.Name = %v, want Pending", dto.Name)
	}
}

func TestFindAllStatusUseCase_MultipleStatuses(t *testing.T) {
	statuses := []struct {
		id   string
		name string
	}{
		{"status-1", "Pending"},
		{"status-2", "Confirmed"},
		{"status-3", "Preparing"},
		{"status-4", "Ready"},
		{"status-5", "Completed"},
	}

	for _, s := range statuses {
		t.Run(s.name, func(t *testing.T) {
			status, err := entities.NewOrderStatus(s.id, s.name)
			if err != nil {
				t.Errorf("NewOrderStatus(%v, %v) unexpected error: %v", s.id, s.name, err)
			}
			if status.ID != s.id {
				t.Errorf("OrderStatus.ID = %v, want %v", status.ID, s.id)
			}
			if status.Name.Value() != s.name {
				t.Errorf("OrderStatus.Name = %v, want %v", status.Name.Value(), s.name)
			}
		})
	}
}

func TestFindAllStatusUseCase_InvalidStatusName(t *testing.T) {
	invalidNames := []string{
		"",
		"ab",
	}

	for _, name := range invalidNames {
		t.Run(name, func(t *testing.T) {
			_, err := entities.NewOrderStatus("status-1", name)
			if err == nil {
				t.Errorf("NewOrderStatus() with name '%v' expected error, got nil", name)
			}
		})
	}
}

func TestNewFindAllOrderStatusUseCase(t *testing.T) {
	_ = NewFindAllOrderStatusUseCase
}
// Comprehensive tests using mocks for full coverage

func TestFindAllOrderStatusUseCase_NewFindAllOrderStatusUseCase(t *testing.T) {
	mockGateway := NewMockOrderStatusGateway()
	uc := NewFindAllOrderStatusUseCase(mockGateway)

	if uc == nil {
		t.Error("Expected use case to be created")
	}
}

func TestFindAllOrderStatusUseCase_Execute_Success(t *testing.T) {
	mockGateway := NewMockOrderStatusGateway()
	uc := NewFindAllOrderStatusUseCase(mockGateway)

	// Add test statuses
	status1, _ := entities.NewOrderStatus("pending", "Pending")
	status2, _ := entities.NewOrderStatus("paid", "Paid")
	
	mockGateway.AddStatus(status1)
	mockGateway.AddStatus(status2)

	statuses, err := uc.Execute()

	if err != nil {
		t.Errorf("Expected no error for successful find all, got %v", err)
	}

	if len(statuses) != 2 {
		t.Errorf("Expected 2 statuses, got %d", len(statuses))
	}
}

func TestFindAllOrderStatusUseCase_Execute_EmptyResult(t *testing.T) {
	mockGateway := NewMockOrderStatusGateway()
	uc := NewFindAllOrderStatusUseCase(mockGateway)

	statuses, err := uc.Execute()

	if err != nil {
		t.Errorf("Expected no error for empty result, got %v", err)
	}

	if len(statuses) != 0 {
		t.Errorf("Expected 0 statuses, got %d", len(statuses))
	}
}