package consumers

import (
	"context"
	"errors"
	"testing"

	"microservice/internal/adapters/brokers"
	"microservice/utils/identity"

	"github.com/stretchr/testify/assert"
)

func TestNewOrderErrorConsumer(t *testing.T) {
	broker := &mockBroker{}
	orderGateway := &mockOrderGateway{}

	consumer := NewOrderErrorConsumer(broker, orderGateway)

	assert.NotNil(t, consumer)
	assert.Equal(t, broker, consumer.broker)
	assert.NotNil(t, consumer.deleteOrderUseCase)
}

func TestOrderErrorConsumer_Start_Success(t *testing.T) {
	broker := &mockBroker{
		consumeOrderErrorFunc: func(ctx context.Context, handler brokers.OrderErrorHandler) error {
			// Simulate successful consumption
			return nil
		},
	}
	orderGateway := &mockOrderGateway{}

	consumer := NewOrderErrorConsumer(broker, orderGateway)

	ctx := context.Background()
	err := consumer.Start(ctx)

	assert.NoError(t, err)
}

func TestOrderErrorConsumer_Start_BrokerError(t *testing.T) {
	broker := &mockBroker{
		consumeOrderErrorFunc: func(ctx context.Context, handler brokers.OrderErrorHandler) error {
			return errors.New("broker connection failed")
		},
	}
	orderGateway := &mockOrderGateway{}

	consumer := NewOrderErrorConsumer(broker, orderGateway)

	ctx := context.Background()
	err := consumer.Start(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "broker connection failed")
}

func TestOrderErrorConsumer_processOrderError_Success(t *testing.T) {
	orderId := identity.NewUUIDV4()

	orderGateway := &mockOrderGateway{
		deleteFunc: func(id string) error {
			if id == orderId {
				return nil
			}
			return errors.New("order not found")
		},
	}

	var capturedHandler brokers.OrderErrorHandler
	broker := &mockBroker{
		consumeOrderErrorFunc: func(ctx context.Context, handler brokers.OrderErrorHandler) error {
			capturedHandler = handler
			return nil
		},
	}

	consumer := NewOrderErrorConsumer(broker, orderGateway)

	ctx := context.Background()
	err := consumer.Start(ctx)
	assert.NoError(t, err)

	// Test the captured handler
	message := brokers.OrderErrorMessage{
		OrderID:         orderId,
		SystemTriggered: "kitchen_order_service",
	}

	err = capturedHandler(message)
	assert.NoError(t, err)
}

func TestOrderErrorConsumer_processOrderError_DeleteError(t *testing.T) {
	orderId := identity.NewUUIDV4()

	orderGateway := &mockOrderGateway{
		deleteFunc: func(id string) error {
			return errors.New("database deletion failed")
		},
	}

	var capturedHandler brokers.OrderErrorHandler
	broker := &mockBroker{
		consumeOrderErrorFunc: func(ctx context.Context, handler brokers.OrderErrorHandler) error {
			capturedHandler = handler
			return nil
		},
	}

	consumer := NewOrderErrorConsumer(broker, orderGateway)

	ctx := context.Background()
	err := consumer.Start(ctx)
	assert.NoError(t, err)

	// Test the captured handler with delete error
	message := brokers.OrderErrorMessage{
		OrderID:         orderId,
		SystemTriggered: "kitchen_order_service",
	}

	err = capturedHandler(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database deletion failed")
}

func TestOrderErrorConsumer_processOrderError_OrderNotFound(t *testing.T) {
	orderId := identity.NewUUIDV4()

	orderGateway := &mockOrderGateway{
		deleteFunc: func(id string) error {
			return errors.New("order not found")
		},
	}

	var capturedHandler brokers.OrderErrorHandler
	broker := &mockBroker{
		consumeOrderErrorFunc: func(ctx context.Context, handler brokers.OrderErrorHandler) error {
			capturedHandler = handler
			return nil
		},
	}

	consumer := NewOrderErrorConsumer(broker, orderGateway)

	ctx := context.Background()
	err := consumer.Start(ctx)
	assert.NoError(t, err)

	// Test the captured handler with non-existent order
	message := brokers.OrderErrorMessage{
		OrderID:         orderId,
		SystemTriggered: "kitchen_order_service",
	}

	err = capturedHandler(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order not found")
}

func TestOrderErrorConsumer_Structure(t *testing.T) {
	consumer := &OrderErrorConsumer{
		broker:             nil,
		deleteOrderUseCase: nil,
	}

	assert.NotNil(t, consumer)
}

func TestOrderErrorMessage_Structure(t *testing.T) {
	orderId := identity.NewUUIDV4()

	message := brokers.OrderErrorMessage{
		OrderID:         orderId,
		SystemTriggered: "kitchen_order_service",
	}

	assert.Equal(t, orderId, message.OrderID)
}
