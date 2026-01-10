package use_cases

import (
	"errors"
	"testing"

	"microservice/internal/adapters/dtos"
	"microservice/internal/domain/entities"

	"github.com/stretchr/testify/assert"
)

// Mock implementations for testing
type mockOrderGateway struct {
	findByIDFunc func(id string) (*entities.Order, error)
	updateFunc   func(order entities.Order) error
	findAllFunc  func(filter dtos.OrderFilterDTO) ([]entities.Order, error)
}

func (m *mockOrderGateway) FindByID(id string) (*entities.Order, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockOrderGateway) Update(order entities.Order) error {
	if m.updateFunc != nil {
		return m.updateFunc(order)
	}
	return nil
}

func (m *mockOrderGateway) Create(order entities.Order) error {
	return errors.New("not implemented")
}

func (m *mockOrderGateway) FindAll(filter dtos.OrderFilterDTO) ([]entities.Order, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(filter)
	}
	return nil, errors.New("not implemented")
}

func (m *mockOrderGateway) Delete(id string) error {
	return errors.New("not implemented")
}

type mockOrderStatusGateway struct {
	findByNameFunc func(name string) (*entities.OrderStatus, error)
	findByIDFunc   func(id string) (*entities.OrderStatus, error)
	findAllFunc    func() ([]entities.OrderStatus, error)
}

func (m *mockOrderStatusGateway) FindByName(name string) (*entities.OrderStatus, error) {
	if m.findByNameFunc != nil {
		return m.findByNameFunc(name)
	}
	return nil, errors.New("not implemented")
}

func (m *mockOrderStatusGateway) FindByID(id string) (*entities.OrderStatus, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockOrderStatusGateway) FindAll() ([]entities.OrderStatus, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc()
	}
	return nil, errors.New("not implemented")
}

func TestNewUpdateOrderStatusUseCase(t *testing.T) {
	orderGateway := &mockOrderGateway{}
	statusGateway := &mockOrderStatusGateway{}
	
	useCase := NewUpdateOrderStatusUseCase(orderGateway, statusGateway)
	assert.NotNil(t, useCase)
	assert.Equal(t, orderGateway, useCase.orderGateway)
	assert.Equal(t, statusGateway, useCase.orderStatusGateway)
}

func TestUpdateOrderStatusUseCase_Execute_Success(t *testing.T) {
	customerID := "customer-123"
	order, _ := entities.NewOrder("order-123", &customerID)
	oldStatus, _ := entities.NewOrderStatus("status-1", "Recebido")
	newStatus, _ := entities.NewOrderStatus("status-2", "Em preparação")
	order.Status = *oldStatus

	orderGateway := &mockOrderGateway{
		findByIDFunc: func(id string) (*entities.Order, error) {
			if id == "order-123" {
				return order, nil
			}
			return nil, errors.New("order not found")
		},
		updateFunc: func(o entities.Order) error {
			return nil
		},
	}

	statusGateway := &mockOrderStatusGateway{
		findByNameFunc: func(name string) (*entities.OrderStatus, error) {
			if name == "Em preparação" {
				return newStatus, nil
			}
			return nil, errors.New("status not found")
		},
	}

	useCase := NewUpdateOrderStatusUseCase(orderGateway, statusGateway)
	
	dto := UpdateOrderStatusDTO{
		OrderID: "order-123",
		Status:  "Em preparação",
	}

	result, err := useCase.Execute(dto)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "order-123", result.Order.ID)
	assert.Equal(t, "Em preparação", result.Order.Status.Name.Value())
	assert.Contains(t, result.Message, "Order order-123 status updated to Em preparação")
}

func TestUpdateOrderStatusUseCase_Execute_OrderNotFound(t *testing.T) {
	orderGateway := &mockOrderGateway{
		findByIDFunc: func(id string) (*entities.Order, error) {
			return nil, errors.New("order not found")
		},
	}

	statusGateway := &mockOrderStatusGateway{}

	useCase := NewUpdateOrderStatusUseCase(orderGateway, statusGateway)
	
	dto := UpdateOrderStatusDTO{
		OrderID: "non-existent-order",
		Status:  "Em preparação",
	}

	result, err := useCase.Execute(dto)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to find order non-existent-order")
}

func TestUpdateOrderStatusUseCase_Execute_StatusNotFound(t *testing.T) {
	customerID := "customer-123"
	order, _ := entities.NewOrder("order-123", &customerID)
	status, _ := entities.NewOrderStatus("status-1", "Recebido")
	order.Status = *status

	orderGateway := &mockOrderGateway{
		findByIDFunc: func(id string) (*entities.Order, error) {
			return order, nil
		},
	}

	statusGateway := &mockOrderStatusGateway{
		findByNameFunc: func(name string) (*entities.OrderStatus, error) {
			return nil, errors.New("status not found")
		},
	}

	useCase := NewUpdateOrderStatusUseCase(orderGateway, statusGateway)
	
	dto := UpdateOrderStatusDTO{
		OrderID: "order-123",
		Status:  "Status Inexistente",
	}

	result, err := useCase.Execute(dto)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to find order status 'Status Inexistente'")
}

func TestUpdateOrderStatusUseCase_Execute_UpdateError(t *testing.T) {
	customerID := "customer-123"
	order, _ := entities.NewOrder("order-123", &customerID)
	oldStatus, _ := entities.NewOrderStatus("status-1", "Recebido")
	newStatus, _ := entities.NewOrderStatus("status-2", "Em preparação")
	order.Status = *oldStatus

	orderGateway := &mockOrderGateway{
		findByIDFunc: func(id string) (*entities.Order, error) {
			return order, nil
		},
		updateFunc: func(o entities.Order) error {
			return errors.New("database update failed")
		},
	}

	statusGateway := &mockOrderStatusGateway{
		findByNameFunc: func(name string) (*entities.OrderStatus, error) {
			return newStatus, nil
		},
	}

	useCase := NewUpdateOrderStatusUseCase(orderGateway, statusGateway)
	
	dto := UpdateOrderStatusDTO{
		OrderID: "order-123",
		Status:  "Em preparação",
	}

	result, err := useCase.Execute(dto)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update order order-123")
}

func TestUpdateOrderStatusUseCase_MapKitchenStatusToOrderStatus(t *testing.T) {
	useCase := NewUpdateOrderStatusUseCase(nil, nil)

	// Test cases para mapeamento de status
	testCases := []struct {
		name           string
		kitchenStatus  string
		expectedStatus string
	}{
		{
			name:           "Em preparação mapping",
			kitchenStatus:  "Em preparação",
			expectedStatus: "Em preparação",
		},
		{
			name:           "Pronto mapping",
			kitchenStatus:  "Pronto",
			expectedStatus: "Pronto",
		},
		{
			name:           "Finalizado mapping",
			kitchenStatus:  "Finalizado",
			expectedStatus: "Finalizado",
		},
		{
			name:           "Unknown status returns same",
			kitchenStatus:  "Status Desconhecido",
			expectedStatus: "Status Desconhecido",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := useCase.mapKitchenStatusToOrderStatus(tc.kitchenStatus)
			assert.Equal(t, tc.expectedStatus, result)
		})
	}
}

func TestUpdateOrderStatusDTO_Structure(t *testing.T) {
	dto := UpdateOrderStatusDTO{
		OrderID: "order-123",
		Status:  "Em preparação",
	}

	assert.Equal(t, "order-123", dto.OrderID)
	assert.Equal(t, "Em preparação", dto.Status)
}

func TestUpdateOrderStatusResult_Structure(t *testing.T) {
	customerID := "customer-123"
	order, _ := entities.NewOrder("order-123", &customerID)
	status, _ := entities.NewOrderStatus("1", "Recebido")
	order.Status = *status

	result := UpdateOrderStatusResult{
		Order:   *order,
		Message: "Order updated successfully",
	}

	assert.Equal(t, "order-123", result.Order.ID)
	assert.Equal(t, "Order updated successfully", result.Message)
}

func TestUpdateOrderStatusUseCase_Execute_AllKitchenStatuses(t *testing.T) {
	testCases := []struct {
		name          string
		kitchenStatus string
		expectedStatus string
	}{
		{"Em preparação", "Em preparação", "Em preparação"},
		{"Pronto", "Pronto", "Pronto"},
		{"Finalizado", "Finalizado", "Finalizado"},
		{"Cancelado", "Cancelado", "Cancelado"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			customerID := "customer-123"
			order, _ := entities.NewOrder("order-123", &customerID)
			oldStatus, _ := entities.NewOrderStatus("status-1", "Recebido")
			newStatus, _ := entities.NewOrderStatus("status-2", tc.expectedStatus)
			order.Status = *oldStatus

			orderGateway := &mockOrderGateway{
				findByIDFunc: func(id string) (*entities.Order, error) {
					return order, nil
				},
				updateFunc: func(o entities.Order) error {
					return nil
				},
			}

			statusGateway := &mockOrderStatusGateway{
				findByNameFunc: func(name string) (*entities.OrderStatus, error) {
					if name == tc.expectedStatus {
						return newStatus, nil
					}
					return nil, errors.New("status not found")
				},
			}

			useCase := NewUpdateOrderStatusUseCase(orderGateway, statusGateway)
			
			dto := UpdateOrderStatusDTO{
				OrderID: "order-123",
				Status:  tc.kitchenStatus,
			}

			result, err := useCase.Execute(dto)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tc.expectedStatus, result.Order.Status.Name.Value())
		})
	}
}