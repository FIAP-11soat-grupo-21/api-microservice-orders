package consumers

import (
	"context"
	"testing"

	"microservice/internal/adapters/brokers"
)

type mockBroker struct{}

func (m *mockBroker) ConsumePaymentConfirmations(ctx context.Context, handler brokers.PaymentConfirmationHandler) error {
	return nil
}

func (m *mockBroker) SendToKitchen(message map[string]interface{}) error {
	return nil
}

func (m *mockBroker) Close() error {
	return nil
}

func TestPaymentConsumerStruct(t *testing.T) {
	broker := &mockBroker{}
	kitchenBroker := &mockBroker{}

	consumer := &PaymentConsumer{
		broker:        broker,
		kitchenBroker: kitchenBroker,
	}

	if consumer == nil {
		t.Error("Expected consumer to be created")
	}

	if consumer.broker == nil {
		t.Error("Expected broker to be set")
	}

	if consumer.kitchenBroker == nil {
		t.Error("Expected kitchenBroker to be set")
	}
}