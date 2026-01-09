package use_cases

import (
	"context"
	"errors"
	"testing"

	"microservice/internal/adapters/brokers"
	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
	"microservice/internal/adapters/gateways"
)

type testOrderDataSource struct {
	orders map[string]daos.OrderDAO
}

func newTestOrderDataSource() *testOrderDataSource {
	return &testOrderDataSource{
		orders: make(map[string]daos.OrderDAO),
	}
}

func (ds *testOrderDataSource) Create(order daos.OrderDAO) error {
	ds.orders[order.ID] = order
	return nil
}

func (ds *testOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	result := make([]daos.OrderDAO, 0, len(ds.orders))
	for _, order := range ds.orders {
		result = append(result, order)
	}
	return result, nil
}

func (ds *testOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	order, ok := ds.orders[id]
	if !ok {
		return daos.OrderDAO{}, errors.New("not found")
	}
	return order, nil
}

func (ds *testOrderDataSource) Update(order daos.OrderDAO) error {
	ds.orders[order.ID] = order
	return nil
}

func (ds *testOrderDataSource) Delete(id string) error {
	delete(ds.orders, id)
	return nil
}

type testOrderStatusDataSource struct {
	statuses map[string]daos.OrderStatusDAO
}

func newTestOrderStatusDataSource() *testOrderStatusDataSource {
	ds := &testOrderStatusDataSource{
		statuses: make(map[string]daos.OrderStatusDAO),
	}
	ds.statuses["56d3b3c3-1801-49cd-bae7-972c78082012"] = daos.OrderStatusDAO{
		ID:   "56d3b3c3-1801-49cd-bae7-972c78082012",
		Name: "Recebido",
	}
	return ds
}

func (ds *testOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	status, ok := ds.statuses[id]
	if !ok {
		return daos.OrderStatusDAO{}, errors.New("not found")
	}
	return status, nil
}

func (ds *testOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	result := make([]daos.OrderStatusDAO, 0, len(ds.statuses))
	for _, status := range ds.statuses {
		result = append(result, status)
	}
	return result, nil
}

// Mock do MessageBroker para testes
type testMessageBroker struct{}

func (m *testMessageBroker) SendToKitchen(message map[string]interface{}) error {
	return nil
}

func (m *testMessageBroker) ConsumePaymentConfirmations(ctx context.Context, handler brokers.PaymentConfirmationHandler) error {
	return nil
}

func (m *testMessageBroker) Close() error {
	return nil
}

func TestCreateOrderUseCase_Execute_Integration(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := newTestOrderStatusDataSource()
	broker := &testMessageBroker{}

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	uc := NewCreateOrderUseCase(*orderGateway, *statusGateway, broker)

	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{ProductID: "product-1", Quantity: 2, Price: 10.0},
		{ProductID: "product-2", Quantity: 1, Price: 25.0},
	}

	result, err := uc.Execute(&customerID, items)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID == "" {
		t.Error("Expected order ID to be set")
	}

	if *result.CustomerID != customerID {
		t.Errorf("Expected customer ID %s, got %s", customerID, *result.CustomerID)
	}

	if len(result.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result.Items))
	}

	expectedAmount := 45.0
	if result.Amount.Value() != expectedAmount {
		t.Errorf("Expected amount %.2f, got %.2f", expectedAmount, result.Amount.Value())
	}
}

func TestCreateOrderUseCase_Execute_WithoutCustomerID_Integration(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := newTestOrderStatusDataSource()
	broker := &testMessageBroker{}

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	uc := NewCreateOrderUseCase(*orderGateway, *statusGateway, broker)

	items := []dtos.CreateOrderItemDTO{
		{ProductID: "product-1", Quantity: 1, Price: 15.0},
	}

	result, err := uc.Execute(nil, items)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID == "" {
		t.Error("Expected order ID to be set")
	}

	if result.CustomerID != nil {
		t.Error("Expected customer ID to be nil")
	}
}

func TestCreateOrderUseCase_Execute_StatusNotFound_Integration(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := &testOrderStatusDataSource{statuses: make(map[string]daos.OrderStatusDAO)} // Empty statuses
	broker := &testMessageBroker{}

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	uc := NewCreateOrderUseCase(*orderGateway, *statusGateway, broker)

	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{ProductID: "product-1", Quantity: 1, Price: 15.0},
	}

	_, err := uc.Execute(&customerID, items)

	if err == nil {
		t.Error("Expected error when status not found")
	}
}

func TestFindAllOrdersUseCase_Execute_Integration(t *testing.T) {
	orderDS := newTestOrderDataSource()

	orderDS.orders["550e8400-e29b-41d4-a716-446655440000"] = daos.OrderDAO{
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
	}

	orderGateway := gateways.NewOrderGateway(orderDS)

	uc := NewFindAllOrdersUseCase(*orderGateway)

	result, err := uc.Execute(dtos.OrderFilterDTO{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 order, got %d", len(result))
	}
}

func TestFindAllOrderStatusUseCase_Execute_Integration(t *testing.T) {
	statusDS := newTestOrderStatusDataSource()

	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	uc := NewFindAllOrderStatusUseCase(*statusGateway)

	result, err := uc.Execute()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected at least one status")
	}
}

// Error data source for testing error scenarios
type errorOrderDataSource struct{}

func (ds *errorOrderDataSource) Create(order daos.OrderDAO) error {
	return errors.New("database error")
}

func (ds *errorOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	return nil, errors.New("database error")
}

func (ds *errorOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	return daos.OrderDAO{}, errors.New("database error")
}

func (ds *errorOrderDataSource) Update(order daos.OrderDAO) error {
	return errors.New("database error")
}

func (ds *errorOrderDataSource) Delete(id string) error {
	return errors.New("database error")
}

func TestCreateOrderUseCase_Execute_DatabaseError_Integration(t *testing.T) {
	orderDS := &errorOrderDataSource{}
	statusDS := newTestOrderStatusDataSource()
	broker := &testMessageBroker{}

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	uc := NewCreateOrderUseCase(*orderGateway, *statusGateway, broker)

	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{ProductID: "product-1", Quantity: 1, Price: 15.0},
	}

	_, err := uc.Execute(&customerID, items)

	if err == nil {
		t.Error("Expected error from database")
	}
}

func TestFindAllOrdersUseCase_Execute_DatabaseError_Integration(t *testing.T) {
	orderDS := &errorOrderDataSource{}

	orderGateway := gateways.NewOrderGateway(orderDS)

	uc := NewFindAllOrdersUseCase(*orderGateway)

	_, err := uc.Execute(dtos.OrderFilterDTO{})

	if err == nil {
		t.Error("Expected error from database")
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
