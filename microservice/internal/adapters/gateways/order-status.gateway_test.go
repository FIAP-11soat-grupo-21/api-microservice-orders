package gateways

import (
	"errors"
	"testing"

	"microservice/internal/adapters/daos"
)

type mockOrderStatusDataSource struct {
	findByIDFunc func(id string) (daos.OrderStatusDAO, error)
	findAllFunc  func() ([]daos.OrderStatusDAO, error)
}

func (m *mockOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return daos.OrderStatusDAO{}, nil
}

func (m *mockOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc()
	}
	return []daos.OrderStatusDAO{}, nil
}

func TestNewOrderStatusGateway(t *testing.T) {
	ds := &mockOrderStatusDataSource{}
	gateway := NewOrderStatusGateway(ds)

	if gateway == nil {
		t.Error("NewOrderStatusGateway() returned nil")
	}
}

func TestOrderStatusGateway_FindAll_Success(t *testing.T) {
	ds := &mockOrderStatusDataSource{
		findAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return []daos.OrderStatusDAO{
				{ID: "status-1", Name: "Pending"},
				{ID: "status-2", Name: "Confirmed"},
			}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	statuses, err := gateway.FindAll()

	if err != nil {
		t.Errorf("FindAll() unexpected error: %v", err)
	}
	if len(statuses) != 2 {
		t.Errorf("FindAll() length = %v, want 2", len(statuses))
	}
	if statuses[0].ID != "status-1" {
		t.Errorf("FindAll()[0].ID = %v, want status-1", statuses[0].ID)
	}
	if statuses[0].Name.Value() != "Pending" {
		t.Errorf("FindAll()[0].Name = %v, want Pending", statuses[0].Name.Value())
	}
}

func TestOrderStatusGateway_FindAll_Error(t *testing.T) {
	ds := &mockOrderStatusDataSource{
		findAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return nil, errors.New("database error")
		},
	}

	gateway := NewOrderStatusGateway(ds)
	_, err := gateway.FindAll()

	if err == nil {
		t.Error("FindAll() expected error, got nil")
	}
}

func TestOrderStatusGateway_FindAll_InvalidName(t *testing.T) {
	ds := &mockOrderStatusDataSource{
		findAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return []daos.OrderStatusDAO{
				{ID: "status-1", Name: "ab"}, // Invalid: too short
			}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	_, err := gateway.FindAll()

	if err == nil {
		t.Error("FindAll() expected error for invalid name, got nil")
	}
}

func TestOrderStatusGateway_FindByID_Success(t *testing.T) {
	ds := &mockOrderStatusDataSource{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{ID: "status-1", Name: "Pending"}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	status, err := gateway.FindByID("status-1")

	if err != nil {
		t.Errorf("FindByID() unexpected error: %v", err)
	}
	if status == nil {
		t.Fatal("FindByID() returned nil")
	}
	if status.ID != "status-1" {
		t.Errorf("FindByID() ID = %v, want status-1", status.ID)
	}
	if status.Name.Value() != "Pending" {
		t.Errorf("FindByID() Name = %v, want Pending", status.Name.Value())
	}
}

func TestOrderStatusGateway_FindByID_Error(t *testing.T) {
	ds := &mockOrderStatusDataSource{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{}, errors.New("not found")
		},
	}

	gateway := NewOrderStatusGateway(ds)
	_, err := gateway.FindByID("status-1")

	if err == nil {
		t.Error("FindByID() expected error, got nil")
	}
}

func TestOrderStatusGateway_FindByID_InvalidName(t *testing.T) {
	ds := &mockOrderStatusDataSource{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{ID: "status-1", Name: "ab"}, nil // Invalid name
		},
	}

	gateway := NewOrderStatusGateway(ds)
	_, err := gateway.FindByID("status-1")

	if err == nil {
		t.Error("FindByID() expected error for invalid name, got nil")
	}
}
