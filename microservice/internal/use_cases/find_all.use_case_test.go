package use_cases

import (
	"testing"
	"time"

	"microservice/internal/adapters/dtos"
)

func TestFindAllOrdersUseCase_FilterDTOStructure(t *testing.T) {
	customerID := "customer-123"
	statusID := "status-1"
	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()

	filter := dtos.OrderFilterDTO{
		CustomerID:    &customerID,
		StatusID:      &statusID,
		CreatedAtFrom: &from,
		CreatedAtTo:   &to,
	}

	if *filter.CustomerID != customerID {
		t.Errorf("OrderFilterDTO.CustomerID = %v, want %v", *filter.CustomerID, customerID)
	}
	if *filter.StatusID != statusID {
		t.Errorf("OrderFilterDTO.StatusID = %v, want %v", *filter.StatusID, statusID)
	}
	if filter.CreatedAtFrom == nil {
		t.Error("OrderFilterDTO.CreatedAtFrom should not be nil")
	}
	if filter.CreatedAtTo == nil {
		t.Error("OrderFilterDTO.CreatedAtTo should not be nil")
	}
}

func TestFindAllOrdersUseCase_EmptyFilter(t *testing.T) {
	filter := dtos.OrderFilterDTO{}

	if filter.CustomerID != nil {
		t.Errorf("Empty OrderFilterDTO.CustomerID = %v, want nil", filter.CustomerID)
	}
	if filter.StatusID != nil {
		t.Errorf("Empty OrderFilterDTO.StatusID = %v, want nil", filter.StatusID)
	}
	if filter.CreatedAtFrom != nil {
		t.Errorf("Empty OrderFilterDTO.CreatedAtFrom = %v, want nil", filter.CreatedAtFrom)
	}
	if filter.CreatedAtTo != nil {
		t.Errorf("Empty OrderFilterDTO.CreatedAtTo = %v, want nil", filter.CreatedAtTo)
	}
}

func TestFindAllOrdersUseCase_PartialFilter(t *testing.T) {
	statusID := "status-1"

	filter := dtos.OrderFilterDTO{
		StatusID: &statusID,
	}

	if filter.CustomerID != nil {
		t.Errorf("Partial OrderFilterDTO.CustomerID = %v, want nil", filter.CustomerID)
	}
	if *filter.StatusID != statusID {
		t.Errorf("Partial OrderFilterDTO.StatusID = %v, want %v", *filter.StatusID, statusID)
	}
}

func TestNewFindAllOrdersUseCase(t *testing.T) {
	_ = NewFindAllOrdersUseCase
}
