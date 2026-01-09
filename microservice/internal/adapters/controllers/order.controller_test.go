package controllers

import (
	"context"
	"errors"
	"testing"

	"microservice/internal/adapters/brokers"
	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
)

type mockOrderDataSource struct {
	createFunc   func(order daos.OrderDAO) error
	findAllFunc  func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error)
	findByIDFunc func(id string) (daos.OrderDAO, error)
	updateFunc   func(order daos.OrderDAO) error
	deleteFunc   func(id string) error
}

func (m *mockOrderDataSource) Create(order daos.OrderDAO) error {
	if m.createFunc != nil {
		return m.createFunc(order)
	}
	return nil
}

func (m *mockOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(filter)
	}
	return []daos.OrderDAO{}, nil
}

func (m *mockOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return daos.OrderDAO{}, nil
}

func (m *mockOrderDataSource) Update(order daos.OrderDAO) error {
	if m.updateFunc != nil {
		return m.updateFunc(order)
	}
	return nil
}

func (m *mockOrderDataSource) Delete(id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

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

type mockMessageBroker struct{}

func (m *mockMessageBroker) SendToKitchen(message map[string]interface{}) error {
	return nil
}

func (m *mockMessageBroker) ConsumePaymentConfirmations(ctx context.Context, handler brokers.PaymentConfirmationHandler) error {
	return nil
}

func (m *mockMessageBroker) Close() error {
	return nil
}

func TestNewOrderController(t *testing.T) {
	orderDS := &mockOrderDataSource{}
	statusDS := &mockOrderStatusDataSource{}
	broker := &mockMessageBroker{}

	controller := NewOrderController(orderDS, statusDS, broker)

	if controller == nil {
		t.Error("NewOrderController() returned nil")
	}
}

func TestOrderController_Create_Success(t *testing.T) {
	orderDS := &mockOrderDataSource{
		createFunc: func(order daos.OrderDAO) error {
			return nil
		},
	}
	statusDS := &mockOrderStatusDataSource{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{ID: "status-1", Name: "Pending"}, nil
		},
	}
	broker := &mockMessageBroker{}

	controller := NewOrderController(orderDS, statusDS, broker)

	dto := dtos.CreateOrderDTO{
		CustomerID: stringPtr("customer-1"),
		Items: []dtos.CreateOrderItemDTO{
			{ProductID: "product-1", Quantity: 2, Price: 10.0},
		},
	}

	result, err := controller.Create(dto)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID == "" {
		t.Error("Expected order ID to be set")
	}
}

func TestOrderController_Create_StatusNotFound(t *testing.T) {
	orderDS := &mockOrderDataSource{}
	statusDS := &mockOrderStatusDataSource{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{}, errors.New("not found")
		},
	}
	broker := &mockMessageBroker{}

	controller := NewOrderController(orderDS, statusDS, broker)

	dto := dtos.CreateOrderDTO{
		CustomerID: stringPtr("customer-1"),
		Items: []dtos.CreateOrderItemDTO{
			{ProductID: "product-1", Quantity: 2, Price: 10.0},
		},
	}

	_, err := controller.Create(dto)

	if err == nil {
		t.Error("Expected error when status not found")
	}
}

func TestOrderController_FindAll_Success(t *testing.T) {
	orderDS := &mockOrderDataSource{
		findAllFunc: func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
			return []daos.OrderDAO{
				{
					ID:         "550e8400-e29b-41d4-a716-446655440000",
					CustomerID: stringPtr("customer-1"),
					Amount:     25.0,
					Status: daos.OrderStatusDAO{
						ID:   "status-1",
						Name: "Pending",
					},
					Items: []daos.OrderItemDAO{
						{
							ID:        "item-1",
							OrderID:   "550e8400-e29b-41d4-a716-446655440000",
							ProductID: "product-1",
							Quantity:  1,
							UnitPrice: 25.0,
						},
					},
				},
			}, nil
		},
	}
	statusDS := &mockOrderStatusDataSource{}
	broker := &mockMessageBroker{}

	controller := NewOrderController(orderDS, statusDS, broker)

	result, err := controller.FindAll(dtos.OrderFilterDTO{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 order, got %d", len(result))
	}
}

func TestOrderController_FindAll_Error(t *testing.T) {
	orderDS := &mockOrderDataSource{
		findAllFunc: func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
			return nil, errors.New("database error")
		},
	}
	statusDS := &mockOrderStatusDataSource{}
	broker := &mockMessageBroker{}

	controller := NewOrderController(orderDS, statusDS, broker)

	_, err := controller.FindAll(dtos.OrderFilterDTO{})

	if err == nil {
		t.Error("Expected error from database")
	}
}

func TestOrderController_FindByID_Success(t *testing.T) {
	orderDS := &mockOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:         "550e8400-e29b-41d4-a716-446655440000",
				CustomerID: stringPtr("customer-1"),
				Amount:     25.0,
				Status: daos.OrderStatusDAO{
					ID:   "status-1",
					Name: "Pending",
				},
				Items: []daos.OrderItemDAO{
					{
						ID:        "item-1",
						OrderID:   "550e8400-e29b-41d4-a716-446655440000",
						ProductID: "product-1",
						Quantity:  1,
						UnitPrice: 25.0,
					},
				},
			}, nil
		},
	}
	statusDS := &mockOrderStatusDataSource{}
	broker := &mockMessageBroker{}

	controller := NewOrderController(orderDS, statusDS, broker)

	result, err := controller.FindByID("550e8400-e29b-41d4-a716-446655440000")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected order ID '550e8400-e29b-41d4-a716-446655440000', got '%s'", result.ID)
	}
}

func TestOrderController_FindByID_NotFound(t *testing.T) {
	orderDS := &mockOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{}, errors.New("not found")
		},
	}
	statusDS := &mockOrderStatusDataSource{}
	broker := &mockMessageBroker{}

	controller := NewOrderController(orderDS, statusDS, broker)

	_, err := controller.FindByID("nonexistent")

	if err == nil {
		t.Error("Expected error when order not found")
	}
}

func TestOrderController_Update_Success(t *testing.T) {
	orderDS := &mockOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:         "550e8400-e29b-41d4-a716-446655440000",
				CustomerID: stringPtr("customer-1"),
				Amount:     25.0,
				Status: daos.OrderStatusDAO{
					ID:   "status-1",
					Name: "Pending",
				},
				Items: []daos.OrderItemDAO{
					{
						ID:        "item-1",
						OrderID:   "550e8400-e29b-41d4-a716-446655440000",
						ProductID: "product-1",
						Quantity:  1,
						UnitPrice: 25.0,
					},
				},
			}, nil
		},
		updateFunc: func(order daos.OrderDAO) error {
			return nil
		},
	}
	statusDS := &mockOrderStatusDataSource{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{ID: "status-2", Name: "Confirmed"}, nil
		},
	}
	broker := &mockMessageBroker{}

	controller := NewOrderController(orderDS, statusDS, broker)

	dto := dtos.UpdateOrderDTO{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		StatusID: "status-2",
	}

	result, err := controller.Update(dto)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected order ID '550e8400-e29b-41d4-a716-446655440000', got '%s'", result.ID)
	}
}

func TestOrderController_Delete_Success(t *testing.T) {
	orderDS := &mockOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:         "550e8400-e29b-41d4-a716-446655440000",
				CustomerID: stringPtr("customer-1"),
				Amount:     25.0,
				Status: daos.OrderStatusDAO{
					ID:   "status-1",
					Name: "Pending",
				},
				Items: []daos.OrderItemDAO{
					{
						ID:        "item-1",
						OrderID:   "550e8400-e29b-41d4-a716-446655440000",
						ProductID: "product-1",
						Quantity:  1,
						UnitPrice: 25.0,
					},
				},
			}, nil
		},
		deleteFunc: func(id string) error {
			return nil
		},
	}
	statusDS := &mockOrderStatusDataSource{}
	broker := &mockMessageBroker{}

	controller := NewOrderController(orderDS, statusDS, broker)

	err := controller.Delete("550e8400-e29b-41d4-a716-446655440000")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestOrderController_Delete_Error(t *testing.T) {
	orderDS := &mockOrderDataSource{
		deleteFunc: func(id string) error {
			return errors.New("delete failed")
		},
	}
	statusDS := &mockOrderStatusDataSource{}
	broker := &mockMessageBroker{}

	controller := NewOrderController(orderDS, statusDS, broker)

	err := controller.Delete("order-1")

	if err == nil {
		t.Error("Expected error when delete fails")
	}
}

func TestOrderController_FindAllStatus_Success(t *testing.T) {
	orderDS := &mockOrderDataSource{}
	statusDS := &mockOrderStatusDataSource{
		findAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return []daos.OrderStatusDAO{
				{ID: "status-1", Name: "Pending"},
				{ID: "status-2", Name: "Confirmed"},
			}, nil
		},
	}
	broker := &mockMessageBroker{}

	controller := NewOrderController(orderDS, statusDS, broker)

	result, err := controller.FindAllStatus()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 statuses, got %d", len(result))
	}
}

func TestOrderController_FindAllStatus_Error(t *testing.T) {
	orderDS := &mockOrderDataSource{}
	statusDS := &mockOrderStatusDataSource{
		findAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return nil, errors.New("database error")
		},
	}
	broker := &mockMessageBroker{}

	controller := NewOrderController(orderDS, statusDS, broker)

	_, err := controller.FindAllStatus()

	if err == nil {
		t.Error("Expected error from database")
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}