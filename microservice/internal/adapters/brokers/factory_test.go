package brokers

import (
	"testing"
)

func TestNewFactory(t *testing.T) {
	factory := NewFactory()
	if factory == nil {
		t.Error("Expected factory to be created")
	}
}

func TestFactory_CreateBroker_SQS(t *testing.T) {
	factory := NewFactory()
	config := BrokerConfig{
		Type:                         "sqs",
		SQSUpdateOrderStatusQueueURL: "http://localhost:4566/000000000000/update-order-status-queue",
		SQSOrderErrorQueueURL:        "http://localhost:4566/000000000000/order-error-queue",
		AWSRegion:                    "us-east-1",
	}

	broker, err := factory.CreateBroker(config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if broker == nil {
		t.Error("Expected broker to be created")
	}

	// Check if it's the correct type
	if _, ok := broker.(*SQSBroker); !ok {
		t.Error("Expected SQSBroker type")
	}
}
func TestFactory_CreateBroker_CaseInsensitive(t *testing.T) {
	factory := NewFactory()

	testCases := []string{"SQS", "sqs", "Sqs", "SqS"}

	for _, brokerType := range testCases {
		config := BrokerConfig{
			Type:                         "sqs",
			SQSUpdateOrderStatusQueueURL: "http://localhost:4566/000000000000/update-order-status-queue",
			SQSOrderErrorQueueURL:        "http://localhost:4566/000000000000/order-error-queue",
			AWSRegion:                    "us-east-1",
		}

		broker, err := factory.CreateBroker(config)
		if err != nil {
			t.Errorf("Expected no error for type '%s', got %v", brokerType, err)
		}

		if broker == nil {
			t.Errorf("Expected broker to be created for type '%s'", brokerType)
		}
	}
}

func TestFactory_CreateBroker_UnsupportedType(t *testing.T) {
	factory := NewFactory()
	config := BrokerConfig{
		Type: "unsupported",
	}

	broker, err := factory.CreateBroker(config)
	if err == nil {
		t.Error("Expected error for unsupported broker type")
	}

	if broker != nil {
		t.Error("Expected broker to be nil for unsupported type")
	}

	expectedError := "unsupported broker type: unsupported"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}
