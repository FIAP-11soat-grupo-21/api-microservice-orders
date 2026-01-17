package controllers

import (
	"context"
	"errors"
	"microservice/internal/adapters/brokers"
	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockOrderDataSource struct {
	mock.Mock
}

func (m *MockOrderDataSource) Create(order daos.OrderDAO) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	args := m.Called(filter)
	return args.Get(0).([]daos.OrderDAO), args.Error(1)
}

func (m *MockOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	args := m.Called(id)
	return args.Get(0).(daos.OrderDAO), args.Error(1)
}

func (m *MockOrderDataSource) Update(order daos.OrderDAO) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderDataSource) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockOrderStatusDataSource struct {
	mock.Mock
}

func (m *MockOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	args := m.Called(id)
	return args.Get(0).(daos.OrderStatusDAO), args.Error(1)
}

func (m *MockOrderStatusDataSource) FindByName(name string) (daos.OrderStatusDAO, error) {
	args := m.Called(name)
	return args.Get(0).(daos.OrderStatusDAO), args.Error(1)
}

func (m *MockOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	args := m.Called()
	return args.Get(0).([]daos.OrderStatusDAO), args.Error(1)
}

type MockMessageBroker struct {
	mock.Mock
}

func (m *MockMessageBroker) ConsumeOrderUpdates(ctx context.Context, handler brokers.OrderUpdateHandler) error {
	args := m.Called(ctx, handler)
	return args.Error(0)
}

func (m *MockMessageBroker) ConsumeOrderError(ctx context.Context, handler brokers.OrderErrorHandler) error {
	args := m.Called(ctx, handler)
	return args.Error(0)
}

func (m *MockMessageBroker) PublishOnTopic(ctx context.Context, topic string, message interface{}) error {
	args := m.Called(ctx, topic, message)
	return args.Error(0)
}

func (m *MockMessageBroker) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewOrderController(t *testing.T) {
	mockOrderDS := &MockOrderDataSource{}
	mockOrderStatusDS := &MockOrderStatusDataSource{}
	mockBroker := &MockMessageBroker{}

	controller := NewOrderController(mockOrderDS, mockOrderStatusDS, mockBroker)

	assert.NotNil(t, controller)
	assert.Equal(t, mockOrderDS, controller.orderDataSource)
	assert.Equal(t, mockOrderStatusDS, controller.orderStatusDataSource)
	assert.Equal(t, mockBroker, controller.messageBroker)
	assert.NotNil(t, controller.orderGateway)
	assert.NotNil(t, controller.orderStatusGateway)
}

func TestOrderController_Create_Success(t *testing.T) {
	mockOrderDS := &MockOrderDataSource{}
	mockOrderStatusDS := &MockOrderStatusDataSource{}
	mockBroker := &MockMessageBroker{}

	controller := NewOrderController(mockOrderDS, mockOrderStatusDS, mockBroker)

	customerID := "customer-123"
	createDTO := dtos.CreateOrderDTO{
		CustomerID: &customerID,
		Items: []dtos.CreateOrderItemDTO{
			{
				ProductID: "product-1",
				Quantity:  2,
				Price:     10.50,
			},
		},
	}

	// Mock expectations
	mockOrderStatusDS.On("FindByID", "56d3b3c3-1801-49cd-bae7-972c78082012").Return(daos.OrderStatusDAO{
		ID:   "56d3b3c3-1801-49cd-bae7-972c78082012",
		Name: "PENDING",
	}, nil)

	mockOrderDS.On("Create", mock.AnythingOfType("daos.OrderDAO")).Return(nil)

	result, err := controller.Create(createDTO)

	assert.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, customerID, *result.CustomerID)
	assert.Equal(t, float64(21), result.Amount) // 2 * 10.50
	assert.Equal(t, "PENDING", result.Status.Name)
	assert.Len(t, result.Items, 1)

	mockOrderDS.AssertExpectations(t)
	mockOrderStatusDS.AssertExpectations(t)
}

func TestOrderController_Create_Error(t *testing.T) {
	mockOrderDS := &MockOrderDataSource{}
	mockOrderStatusDS := &MockOrderStatusDataSource{}
	mockBroker := &MockMessageBroker{}

	controller := NewOrderController(mockOrderDS, mockOrderStatusDS, mockBroker)

	customerID := "customer-123"
	createDTO := dtos.CreateOrderDTO{
		CustomerID: &customerID,
		Items: []dtos.CreateOrderItemDTO{
			{
				ProductID: "product-1",
				Quantity:  2,
				Price:     10.50,
			},
		},
	}

	// Mock expectations - simulate error
	mockOrderStatusDS.On("FindByID", "56d3b3c3-1801-49cd-bae7-972c78082012").Return(daos.OrderStatusDAO{}, errors.New("status not found"))

	result, err := controller.Create(createDTO)

	assert.Error(t, err)
	assert.Equal(t, dtos.OrderResponseDTO{}, result)
	assert.Contains(t, err.Error(), "Order Status not found")

	mockOrderStatusDS.AssertExpectations(t)
}

func TestOrderController_FindAll_Success(t *testing.T) {
	mockOrderDS := &MockOrderDataSource{}
	mockOrderStatusDS := &MockOrderStatusDataSource{}
	mockBroker := &MockMessageBroker{}

	controller := NewOrderController(mockOrderDS, mockOrderStatusDS, mockBroker)

	filter := dtos.OrderFilterDTO{}
	now := time.Now()

	mockOrders := []daos.OrderDAO{
		{
			ID:         "order-1",
			CustomerID: stringPtr("customer-1"),
			Amount:     25.50,
			Status: daos.OrderStatusDAO{
				ID:   "status-1",
				Name: "PENDING",
			},
			CreatedAt: now,
			Items: []daos.OrderItemDAO{
				{
					ID:        "item-1",
					ProductID: "product-1",
					OrderID:   "order-1",
					Quantity:  2,
					UnitPrice: 12.75,
				},
			},
		},
	}

	mockOrderDS.On("FindAll", filter).Return(mockOrders, nil)

	result, err := controller.FindAll(filter)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "order-1", result[0].ID)
	assert.Equal(t, "customer-1", *result[0].CustomerID)
	assert.Equal(t, float64(25.50), result[0].Amount)
	assert.Equal(t, "PENDING", result[0].Status.Name)

	mockOrderDS.AssertExpectations(t)
}

func TestOrderController_FindAll_Error(t *testing.T) {
	mockOrderDS := &MockOrderDataSource{}
	mockOrderStatusDS := &MockOrderStatusDataSource{}
	mockBroker := &MockMessageBroker{}

	controller := NewOrderController(mockOrderDS, mockOrderStatusDS, mockBroker)

	filter := dtos.OrderFilterDTO{}

	mockOrderDS.On("FindAll", filter).Return([]daos.OrderDAO{}, errors.New("database error"))

	result, err := controller.FindAll(filter)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")

	mockOrderDS.AssertExpectations(t)
}

func TestOrderController_FindByID_Success(t *testing.T) {
	mockOrderDS := &MockOrderDataSource{}
	mockOrderStatusDS := &MockOrderStatusDataSource{}
	mockBroker := &MockMessageBroker{}

	controller := NewOrderController(mockOrderDS, mockOrderStatusDS, mockBroker)

	orderID := "550e8400-e29b-41d4-a716-446655440000" // Valid UUID
	now := time.Now()

	mockOrder := daos.OrderDAO{
		ID:         orderID,
		CustomerID: stringPtr("customer-1"),
		Amount:     15.75,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "PENDING",
		},
		CreatedAt: now,
		Items: []daos.OrderItemDAO{
			{
				ID:        "item-1",
				ProductID: "product-1",
				OrderID:   orderID,
				Quantity:  1,
				UnitPrice: 15.75,
			},
		},
	}

	mockOrderDS.On("FindByID", orderID).Return(mockOrder, nil)

	result, err := controller.FindByID(orderID)

	assert.NoError(t, err)
	assert.Equal(t, orderID, result.ID)
	assert.Equal(t, "customer-1", *result.CustomerID)
	assert.Equal(t, float64(15.75), result.Amount)
	assert.Equal(t, "PENDING", result.Status.Name)

	mockOrderDS.AssertExpectations(t)
}

func TestOrderController_FindByID_Error(t *testing.T) {
	mockOrderDS := &MockOrderDataSource{}
	mockOrderStatusDS := &MockOrderStatusDataSource{}
	mockBroker := &MockMessageBroker{}

	controller := NewOrderController(mockOrderDS, mockOrderStatusDS, mockBroker)

	orderID := "invalid-order-id" // Invalid UUID

	result, err := controller.FindByID(orderID)

	assert.Error(t, err)
	assert.Equal(t, dtos.OrderResponseDTO{}, result)
	assert.Contains(t, err.Error(), "Invalid order ID")
}

func TestOrderController_Update_Success(t *testing.T) {
	mockOrderDS := &MockOrderDataSource{}
	mockOrderStatusDS := &MockOrderStatusDataSource{}
	mockBroker := &MockMessageBroker{}

	controller := NewOrderController(mockOrderDS, mockOrderStatusDS, mockBroker)

	updateDTO := dtos.UpdateOrderDTO{
		ID:       "550e8400-e29b-41d4-a716-446655440000", // Valid UUID
		StatusID: "status-2",
	}

	now := time.Now()
	mockOrder := daos.OrderDAO{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		CustomerID: stringPtr("customer-1"),
		Amount:     20.00,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "PENDING",
		},
		CreatedAt: now,
		Items:     []daos.OrderItemDAO{},
	}

	mockOrderDS.On("FindByID", "550e8400-e29b-41d4-a716-446655440000").Return(mockOrder, nil)
	mockOrderStatusDS.On("FindByID", "status-2").Return(daos.OrderStatusDAO{
		ID:   "status-2",
		Name: "CONFIRMED",
	}, nil)
	mockOrderDS.On("Update", mock.AnythingOfType("daos.OrderDAO")).Return(nil)

	result, err := controller.Update(updateDTO)

	assert.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.ID)
	assert.Equal(t, "CONFIRMED", result.Status.Name)

	mockOrderDS.AssertExpectations(t)
	mockOrderStatusDS.AssertExpectations(t)
}

func TestOrderController_UpdateStatus_Success(t *testing.T) {
	mockOrderDS := &MockOrderDataSource{}
	mockOrderStatusDS := &MockOrderStatusDataSource{}
	mockBroker := &MockMessageBroker{}

	controller := NewOrderController(mockOrderDS, mockOrderStatusDS, mockBroker)

	updateDTO := dtos.UpdateOrderStatusDTO{
		OrderID: "550e8400-e29b-41d4-a716-446655440000", // Valid UUID
		Status:  "CONFIRMED",
	}

	now := time.Now()
	mockOrder := daos.OrderDAO{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		CustomerID: stringPtr("customer-1"),
		Amount:     20.00,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "PENDING",
		},
		CreatedAt: now,
		Items:     []daos.OrderItemDAO{},
	}

	mockOrderDS.On("FindByID", "550e8400-e29b-41d4-a716-446655440000").Return(mockOrder, nil)
	mockOrderStatusDS.On("FindByName", "CONFIRMED").Return(daos.OrderStatusDAO{
		ID:   "status-2",
		Name: "CONFIRMED",
	}, nil)
	mockOrderDS.On("Update", mock.AnythingOfType("daos.OrderDAO")).Return(nil)

	result, err := controller.UpdateStatus(updateDTO)

	assert.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.ID)
	assert.Equal(t, "CONFIRMED", result.Status.Name)

	mockOrderDS.AssertExpectations(t)
	mockOrderStatusDS.AssertExpectations(t)
}

func TestOrderController_Delete_Success(t *testing.T) {
	mockOrderDS := &MockOrderDataSource{}
	mockOrderStatusDS := &MockOrderStatusDataSource{}
	mockBroker := &MockMessageBroker{}

	controller := NewOrderController(mockOrderDS, mockOrderStatusDS, mockBroker)

	orderID := "550e8400-e29b-41d4-a716-446655440000" // Valid UUID
	now := time.Now()

	mockOrder := daos.OrderDAO{
		ID:         orderID,
		CustomerID: stringPtr("customer-1"),
		Amount:     25.50,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "PENDING",
		},
		CreatedAt: now,
		Items: []daos.OrderItemDAO{
			{
				ID:        "item-1",
				ProductID: "product-1",
				OrderID:   orderID,
				Quantity:  2,
				UnitPrice: 12.75,
			},
		},
	}

	mockOrderDS.On("FindByID", orderID).Return(mockOrder, nil)
	mockOrderDS.On("Delete", orderID).Return(nil)

	err := controller.Delete(orderID)

	assert.NoError(t, err)

	mockOrderDS.AssertExpectations(t)
}

func TestOrderController_Delete_Error(t *testing.T) {
	mockOrderDS := &MockOrderDataSource{}
	mockOrderStatusDS := &MockOrderStatusDataSource{}
	mockBroker := &MockMessageBroker{}

	controller := NewOrderController(mockOrderDS, mockOrderStatusDS, mockBroker)

	orderID := "invalid-order-id" // Invalid UUID

	err := controller.Delete(orderID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid order ID")
}

func TestOrderController_FindAllStatus_Success(t *testing.T) {
	mockOrderDS := &MockOrderDataSource{}
	mockOrderStatusDS := &MockOrderStatusDataSource{}
	mockBroker := &MockMessageBroker{}

	controller := NewOrderController(mockOrderDS, mockOrderStatusDS, mockBroker)

	mockStatuses := []daos.OrderStatusDAO{
		{ID: "status-1", Name: "PENDING"},
		{ID: "status-2", Name: "CONFIRMED"},
		{ID: "status-3", Name: "CANCELLED"},
	}

	mockOrderStatusDS.On("FindAll").Return(mockStatuses, nil)

	result, err := controller.FindAllStatus()

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "PENDING", result[0].Name)
	assert.Equal(t, "CONFIRMED", result[1].Name)
	assert.Equal(t, "CANCELLED", result[2].Name)

	mockOrderStatusDS.AssertExpectations(t)
}

func TestOrderController_FindAllStatus_Error(t *testing.T) {
	mockOrderDS := &MockOrderDataSource{}
	mockOrderStatusDS := &MockOrderStatusDataSource{}
	mockBroker := &MockMessageBroker{}

	controller := NewOrderController(mockOrderDS, mockOrderStatusDS, mockBroker)

	mockOrderStatusDS.On("FindAll").Return([]daos.OrderStatusDAO{}, errors.New("database error"))

	result, err := controller.FindAllStatus()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")

	mockOrderStatusDS.AssertExpectations(t)
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
