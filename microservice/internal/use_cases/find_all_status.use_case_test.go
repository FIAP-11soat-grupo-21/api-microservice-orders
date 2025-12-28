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
