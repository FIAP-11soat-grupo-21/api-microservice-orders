package consumers

import (
	"context"
	"errors"
	"testing"

	"microservice/internal/adapters/brokers"
	"microservice/internal/adapters/dtos"
	"microservice/internal/domain/entities"

	"github.com/stretchr/testify/assert"
)

// Mock implementations for testing
type mockBroker struct {
	consumeOrderUpdatesFunc func(ctx context.Context, handler brokers.OrderUpdateHandler) error
	consumeOrderErrorFunc   func(ctx context.Context, handler brokers.OrderErrorHandler) error
	publishOnTopicFunc      func(ctx context.Context, topic string, message interface{}) error
}

func (m *mockBroker) ConsumeOrderUpdates(ctx context.Context, handler brokers.OrderUpdateHandler) error {
	if m.consumeOrderUpdatesFunc != nil {
		return m.consumeOrderUpdatesFunc(ctx, handler)
	}
	return nil
}

func (m *mockBroker) ConsumeOrderError(ctx context.Context, handler brokers.OrderErrorHandler) error {
	if m.consumeOrderErrorFunc != nil {
		return m.consumeOrderErrorFunc(ctx, handler)
	}
	return nil
}

func (m *mockBroker) PublishOnTopic(ctx context.Context, topic string, message interface{}) error {
	if m.publishOnTopicFunc != nil {
		return m.publishOnTopicFunc(ctx, topic, message)
	}
	return nil
}

func (m *mockBroker) Close() error {
	return nil
}

type mockOrderGateway struct {
	findByIDFunc func(id string) (*entities.Order, error)
	updateFunc   func(order entities.Order) error
}

func (m *mockOrderGateway) FindByID(id string) (*entities.Order, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return nil, nil
}

func (m *mockOrderGateway) Update(order entities.Order) error {
	if m.updateFunc != nil {
		return m.updateFunc(order)
	}
	return nil
}

func (m *mockOrderGateway) Create(order entities.Order) error {
	return nil
}

func (m *mockOrderGateway) FindAll(filter dtos.OrderFilterDTO) ([]entities.Order, error) {
	return nil, nil
}

func (m *mockOrderGateway) Delete(id string) error {
	return nil
}

type mockOrderStatusGateway struct {
	findByNameFunc func(name string) (*entities.OrderStatus, error)
}

func (m *mockOrderStatusGateway) FindByName(name string) (*entities.OrderStatus, error) {
	if m.findByNameFunc != nil {
		return m.findByNameFunc(name)
	}
	return nil, nil
}

func (m *mockOrderStatusGateway) FindByID(id string) (*entities.OrderStatus, error) {
	return nil, nil
}

func (m *mockOrderStatusGateway) FindAll() ([]entities.OrderStatus, error) {
	return nil, nil
}

func TestNewOrderUpdatesConsumer(t *testing.T) {
	broker := &mockBroker{}
	orderGateway := &mockOrderGateway{}
	statusGateway := &mockOrderStatusGateway{}

	consumer := NewOrderUpdatesConsumer(broker, orderGateway, statusGateway)

	assert.NotNil(t, consumer)
	assert.Equal(t, broker, consumer.broker)
	assert.NotNil(t, consumer.updateOrderStatusUseCase)
}

func TestOrderUpdatesConsumer_Start_Success(t *testing.T) {
	broker := &mockBroker{
		consumeOrderUpdatesFunc: func(ctx context.Context, handler brokers.OrderUpdateHandler) error {
			// Simulate successful consumption
			return nil
		},
	}
	orderGateway := &mockOrderGateway{}
	statusGateway := &mockOrderStatusGateway{}

	consumer := NewOrderUpdatesConsumer(broker, orderGateway, statusGateway)

	ctx := context.Background()
	err := consumer.Start(ctx)

	assert.NoError(t, err)
}

func TestOrderUpdatesConsumer_Start_BrokerError(t *testing.T) {
	broker := &mockBroker{
		consumeOrderUpdatesFunc: func(ctx context.Context, handler brokers.OrderUpdateHandler) error {
			return errors.New("broker connection failed")
		},
	}
	orderGateway := &mockOrderGateway{}
	statusGateway := &mockOrderStatusGateway{}

	consumer := NewOrderUpdatesConsumer(broker, orderGateway, statusGateway)

	ctx := context.Background()
	err := consumer.Start(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "broker connection failed")
}

func TestOrderUpdatesConsumer_processOrderUpdate_Success(t *testing.T) {
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

	var capturedHandler brokers.OrderUpdateHandler
	broker := &mockBroker{
		consumeOrderUpdatesFunc: func(ctx context.Context, handler brokers.OrderUpdateHandler) error {
			capturedHandler = handler
			return nil
		},
	}

	consumer := NewOrderUpdatesConsumer(broker, orderGateway, statusGateway)

	ctx := context.Background()
	err := consumer.Start(ctx)
	assert.NoError(t, err)

	// Test the captured handler
	message := brokers.OrderUpdateMessage{
		OrderID: "order-123",
		Status:  "Em preparação",
	}

	err = capturedHandler(message)
	assert.NoError(t, err)
}

func TestOrderUpdatesConsumer_processOrderUpdate_OrderNotFound(t *testing.T) {
	orderGateway := &mockOrderGateway{
		findByIDFunc: func(id string) (*entities.Order, error) {
			return nil, errors.New("order not found")
		},
	}

	statusGateway := &mockOrderStatusGateway{}

	var capturedHandler brokers.OrderUpdateHandler
	broker := &mockBroker{
		consumeOrderUpdatesFunc: func(ctx context.Context, handler brokers.OrderUpdateHandler) error {
			capturedHandler = handler
			return nil
		},
	}

	consumer := NewOrderUpdatesConsumer(broker, orderGateway, statusGateway)

	ctx := context.Background()
	err := consumer.Start(ctx)
	assert.NoError(t, err)

	// Test the captured handler with non-existent order
	message := brokers.OrderUpdateMessage{
		OrderID: "non-existent-order",
		Status:  "Em preparação",
	}

	err = capturedHandler(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find order")
}

func TestOrderUpdatesConsumer_processOrderUpdate_StatusNotFound(t *testing.T) {
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

	var capturedHandler brokers.OrderUpdateHandler
	broker := &mockBroker{
		consumeOrderUpdatesFunc: func(ctx context.Context, handler brokers.OrderUpdateHandler) error {
			capturedHandler = handler
			return nil
		},
	}

	consumer := NewOrderUpdatesConsumer(broker, orderGateway, statusGateway)

	ctx := context.Background()
	err := consumer.Start(ctx)
	assert.NoError(t, err)

	// Test the captured handler with invalid status
	message := brokers.OrderUpdateMessage{
		OrderID: "order-123",
		Status:  "Status Inexistente",
	}

	err = capturedHandler(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find order status")
}

func TestOrderUpdatesConsumer_processOrderUpdate_UpdateError(t *testing.T) {
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

	var capturedHandler brokers.OrderUpdateHandler
	broker := &mockBroker{
		consumeOrderUpdatesFunc: func(ctx context.Context, handler brokers.OrderUpdateHandler) error {
			capturedHandler = handler
			return nil
		},
	}

	consumer := NewOrderUpdatesConsumer(broker, orderGateway, statusGateway)

	ctx := context.Background()
	err := consumer.Start(ctx)
	assert.NoError(t, err)

	// Test the captured handler with update error
	message := brokers.OrderUpdateMessage{
		OrderID: "order-123",
		Status:  "Em preparação",
	}

	err = capturedHandler(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update order")
}

func TestOrderUpdatesConsumer_Structure(t *testing.T) {
	consumer := &OrderUpdatesConsumer{
		broker:                   nil,
		updateOrderStatusUseCase: nil,
	}

	assert.NotNil(t, consumer)
}

func TestOrderUpdateMessage_Structure(t *testing.T) {
	message := brokers.OrderUpdateMessage{
		OrderID: "order-123",
		Status:  "Em preparação",
	}

	assert.Equal(t, "order-123", message.OrderID)
	assert.Equal(t, "Em preparação", message.Status)
}
