package gateways

import (
	"errors"
	"testing"

	"microservice/internal/adapters/daos"
)

type mockOrderStatusDataSource struct {
	findAllFunc    func() ([]daos.OrderStatusDAO, error)
	findByIDFunc   func(id string) (daos.OrderStatusDAO, error)
	findByNameFunc func(name string) (daos.OrderStatusDAO, error)
}

func (m *mockOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc()
	}
	return []daos.OrderStatusDAO{}, nil
}

func (m *mockOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return daos.OrderStatusDAO{}, nil
}

func (m *mockOrderStatusDataSource) FindByName(name string) (daos.OrderStatusDAO, error) {
	if m.findByNameFunc != nil {
		return m.findByNameFunc(name)
	}
	return daos.OrderStatusDAO{}, nil
}

func TestNewOrderStatusGateway(t *testing.T) {
	ds := &mockOrderStatusDataSource{}
	gateway := NewOrderStatusGateway(ds)

	if gateway == nil {
		t.Error("NewOrderStatusGateway() returned nil")
	}
}

/*
|--------------------------------------------------------------------------
| FindAll
|--------------------------------------------------------------------------
*/

func TestOrderStatusGateway_FindAll_Success(t *testing.T) {
	ds := &mockOrderStatusDataSource{
		findAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return []daos.OrderStatusDAO{
				{ID: "status-1", Name: "Pending"},
				{ID: "status-2", Name: "Paid"},
			}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	statuses, err := gateway.FindAll()

	if err != nil {
		t.Errorf("FindAll() unexpected error: %v", err)
	}

	if len(statuses) != 2 {
		t.Errorf("FindAll() length = %d, want 2", len(statuses))
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

func TestOrderStatusGateway_FindAll_InvalidStatus(t *testing.T) {
	ds := &mockOrderStatusDataSource{
		findAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return []daos.OrderStatusDAO{
				{ID: "status-1", Name: "ab"}, // inválido
			}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	_, err := gateway.FindAll()

	if err == nil {
		t.Error("FindAll() expected error for invalid status name, got nil")
	}
}

func TestOrderStatusGateway_FindByID_Success(t *testing.T) {
	ds := &mockOrderStatusDataSource{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{
				ID:   "status-1",
				Name: "Pending",
			}, nil
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
		t.Errorf("FindByID() ID = %s, want status-1", status.ID)
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

func TestOrderStatusGateway_FindByID_InvalidStatus(t *testing.T) {
	ds := &mockOrderStatusDataSource{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{
				ID:   "status-1",
				Name: "ab", // inválido
			}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	_, err := gateway.FindByID("status-1")

	if err == nil {
		t.Error("FindByID() expected error for invalid status name, got nil")
	}
}

func TestOrderStatusGateway_FindByName_Success(t *testing.T) {
	ds := &mockOrderStatusDataSource{
		findByNameFunc: func(name string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{
				ID:   "status-1",
				Name: "Pending",
			}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	status, err := gateway.FindByName("Pending")

	if err != nil {
		t.Errorf("FindByName() unexpected error: %v", err)
	}

	if status == nil {
		t.Fatal("FindByName() returned nil")
	}
}

func TestOrderStatusGateway_FindByName_Error(t *testing.T) {
	ds := &mockOrderStatusDataSource{
		findByNameFunc: func(name string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{}, errors.New("not found")
		},
	}

	gateway := NewOrderStatusGateway(ds)
	_, err := gateway.FindByName("Pending")

	if err == nil {
		t.Error("FindByName() expected error, got nil")
	}
}

func TestOrderStatusGateway_FindByName_InvalidStatus(t *testing.T) {
	ds := &mockOrderStatusDataSource{
		findByNameFunc: func(name string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{
				ID:   "status-1",
				Name: "ab", // inválido
			}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	_, err := gateway.FindByName("ab")

	if err == nil {
		t.Error("FindByName() expected error for invalid status name, got nil")
	}
}
