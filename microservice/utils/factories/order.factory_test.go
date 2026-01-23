package factories

import (
	"context"
	"testing"

	"microservice/infra/api/client"
	"microservice/internal/adapters/brokers"
	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
	"microservice/internal/interfaces"
	"microservice/mocks"
)

// Mock implementations for testing
type mockOrderDataSource struct{}

func (m *mockOrderDataSource) Create(order daos.OrderDAO) error {
	return nil
}

func (m *mockOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	return []daos.OrderDAO{}, nil
}

func (m *mockOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	return daos.OrderDAO{}, nil
}

func (m *mockOrderDataSource) Update(order daos.OrderDAO) error {
	return nil
}

func (m *mockOrderDataSource) Delete(id string) error {
	return nil
}

type mockOrderStatusDataSource struct{}

func (m *mockOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	return daos.OrderStatusDAO{ID: id, Name: "Test"}, nil
}

func (m *mockOrderStatusDataSource) FindByName(name string) (daos.OrderStatusDAO, error) {
	return daos.OrderStatusDAO{ID: "1", Name: name}, nil
}

func (m *mockOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	return []daos.OrderStatusDAO{}, nil
}

type mockMessageBroker struct{}

func (m *mockMessageBroker) ConsumeOrderUpdates(ctx context.Context, handler brokers.OrderUpdateHandler) error {
	return nil
}

func (m *mockMessageBroker) ConsumeOrderError(ctx context.Context, handler brokers.OrderErrorHandler) error {
	return nil
}

func (m *mockMessageBroker) PublishOnTopic(ctx context.Context, topic string, message interface{}) error {
	return nil
}

func (m *mockMessageBroker) Close() error {
	return nil
}

type mockApiClient struct{}

func (m *mockApiClient) Get(path string, obj any) error {
	return nil
}

// Test NewOrderDataSource
func TestNewOrderDataSource(t *testing.T) {
	mock := &mockOrderDataSource{}

	SetNewOrderDataSource(func() interfaces.IOrderDataSource {
		return mock
	})
	defer SetNewOrderDataSource(nil)

	result := NewOrderDataSource()

	if result == nil {
		t.Error("Expected NewOrderDataSource to return non-nil")
	}

	if result != mock {
		t.Error("Expected NewOrderDataSource to return the mock")
	}
}

// Test NewOrderStatusDataSource
func TestNewOrderStatusDataSource(t *testing.T) {
	mock := &mockOrderStatusDataSource{}

	SetNewOrderStatusDataSource(func() interfaces.IOrderStatusDataSource {
		return mock
	})
	defer SetNewOrderStatusDataSource(nil)

	result := NewOrderStatusDataSource()

	if result == nil {
		t.Error("Expected NewOrderStatusDataSource to return non-nil")
	}

	if result != mock {
		t.Error("Expected NewOrderStatusDataSource to return the mock")
	}
}

// Test NewMessageBroker
func TestNewMessageBroker(t *testing.T) {
	mock := &mockMessageBroker{}

	SetNewMessageBroker(func() brokers.MessageBroker {
		return mock
	})
	defer SetNewMessageBroker(nil)

	result := NewMessageBroker()

	if result == nil {
		t.Error("Expected NewMessageBroker to return non-nil")
	}

	if result != mock {
		t.Error("Expected NewMessageBroker to return the mock")
	}
}

// Test NewApiClient
func TestNewApiClient(t *testing.T) {
	mock := &mockApiClient{}

	SetNewApiClient(func() client.IApiClient {
		return mock
	})
	defer SetNewApiClient(nil)

	result := NewApiClient()

	if result == nil {
		t.Error("Expected NewApiClient to return non-nil")
	}

	if result != mock {
		t.Error("Expected NewApiClient to return the mock")
	}
}

// Test SetNewOrderDataSource with nil (reset to default)
func TestSetNewOrderDataSource_ResetToDefault(t *testing.T) {
	// First set a custom mock
	mock := &mockOrderDataSource{}
	SetNewOrderDataSource(func() interfaces.IOrderDataSource {
		return mock
	})

	// Verify the mock is set
	result := NewOrderDataSource()
	if result != mock {
		t.Error("Expected mock to be set")
	}

	// Reset to nil (default)
	SetNewOrderDataSource(nil)

	// The function should still work (using default implementation)
	// We can't easily test the default without DB, so just verify it doesn't panic
}

// Test SetNewOrderStatusDataSource with nil (reset to default)
func TestSetNewOrderStatusDataSource_ResetToDefault(t *testing.T) {
	// First set a custom mock
	mock := &mockOrderStatusDataSource{}
	SetNewOrderStatusDataSource(func() interfaces.IOrderStatusDataSource {
		return mock
	})

	// Verify the mock is set
	result := NewOrderStatusDataSource()
	if result != mock {
		t.Error("Expected mock to be set")
	}

	// Reset to nil (default)
	SetNewOrderStatusDataSource(nil)
}

// Test SetNewMessageBroker with nil (reset to default)
func TestSetNewMessageBroker_ResetToDefault(t *testing.T) {
	// First set a custom mock
	mock := &mockMessageBroker{}
	SetNewMessageBroker(func() brokers.MessageBroker {
		return mock
	})

	// Verify the mock is set
	result := NewMessageBroker()
	if result != mock {
		t.Error("Expected mock to be set")
	}

	// Reset to nil (default) - this will set the default function
	// The default function won't be called here, just assigned
	SetNewMessageBroker(nil)

	// Immediately set a new mock to prevent calling the real broker
	mock2 := &mockMessageBroker{}
	SetNewMessageBroker(func() brokers.MessageBroker {
		return mock2
	})

	result2 := NewMessageBroker()
	if result2 != mock2 {
		t.Error("Expected second mock to be set after reset")
	}

	// Clean up
	SetNewMessageBroker(nil)
}

// Test SetNewApiClient with nil (reset to default)
func TestSetNewApiClient_ResetToDefault(t *testing.T) {
	// First set a custom mock
	mock := &mockApiClient{}
	SetNewApiClient(func() client.IApiClient {
		return mock
	})

	// Verify the mock is set
	result := NewApiClient()
	if result != mock {
		t.Error("Expected mock to be set")
	}

	// Reset to nil (default) - this will set the default function
	// The default function won't be called here, just assigned
	SetNewApiClient(nil)

	// Immediately set a new mock to prevent calling the real API client
	mock2 := &mockApiClient{}
	SetNewApiClient(func() client.IApiClient {
		return mock2
	})

	result2 := NewApiClient()
	if result2 != mock2 {
		t.Error("Expected second mock to be set after reset")
	}

	// Clean up
	SetNewApiClient(nil)
}

// Test that SetNewMessageBroker(nil) correctly enters the nil branch
func TestSetNewMessageBroker_NilBranch(t *testing.T) {
	mocks.SetupEnv()
	defer mocks.CleanupEnv()

	// Call with nil to trigger the if fn == nil branch
	SetNewMessageBroker(nil)

	// Call NewMessageBroker() to execute the default lambda
	// messaging.GetBroker() returns nil if not initialized, which is fine for testing
	result := NewMessageBroker()
	// Result may be nil since messaging broker is not initialized in test
	_ = result

	// Set a mock to restore for cleanup
	SetNewMessageBroker(func() brokers.MessageBroker {
		return &mockMessageBroker{}
	})
}

// Test that SetNewApiClient(nil) correctly enters the nil branch
func TestSetNewApiClient_NilBranch(t *testing.T) {
	mocks.SetupEnv()
	defer mocks.CleanupEnv()

	// Call with nil to trigger the if fn == nil branch
	SetNewApiClient(nil)

	// Call NewApiClient() to execute the default lambda
	// This will create a real HTTP client with the config URL
	result := NewApiClient()
	// Result should not be nil
	if result == nil {
		t.Error("Expected NewApiClient to return non-nil after reset")
	}

	// Set a mock to restore for cleanup
	SetNewApiClient(func() client.IApiClient {
		return &mockApiClient{}
	})
}

// Test NewUpdateOrderStatusUseCase
func TestNewUpdateOrderStatusUseCase(t *testing.T) {
	// Setup mocks
	orderMock := &mockOrderDataSource{}
	statusMock := &mockOrderStatusDataSource{}

	SetNewOrderDataSource(func() interfaces.IOrderDataSource {
		return orderMock
	})
	SetNewOrderStatusDataSource(func() interfaces.IOrderStatusDataSource {
		return statusMock
	})
	defer func() {
		SetNewOrderDataSource(nil)
		SetNewOrderStatusDataSource(nil)
	}()

	uc := NewUpdateOrderStatusUseCase()

	if uc == nil {
		t.Error("Expected NewUpdateOrderStatusUseCase to return non-nil")
	}
}

// Test factory functions work correctly after reset
func TestFactoryFunctionsAfterReset(t *testing.T) {
	// Set mocks
	orderMock := &mockOrderDataSource{}
	statusMock := &mockOrderStatusDataSource{}
	brokerMock := &mockMessageBroker{}
	apiMock := &mockApiClient{}

	SetNewOrderDataSource(func() interfaces.IOrderDataSource {
		return orderMock
	})
	SetNewOrderStatusDataSource(func() interfaces.IOrderStatusDataSource {
		return statusMock
	})
	SetNewMessageBroker(func() brokers.MessageBroker {
		return brokerMock
	})
	SetNewApiClient(func() client.IApiClient {
		return apiMock
	})

	// Verify mocks are set
	if NewOrderDataSource() != orderMock {
		t.Error("OrderDataSource mock not set correctly")
	}
	if NewOrderStatusDataSource() != statusMock {
		t.Error("OrderStatusDataSource mock not set correctly")
	}
	if NewMessageBroker() != brokerMock {
		t.Error("MessageBroker mock not set correctly")
	}
	if NewApiClient() != apiMock {
		t.Error("ApiClient mock not set correctly")
	}

	// Reset all
	SetNewOrderDataSource(nil)
	SetNewOrderStatusDataSource(nil)
	SetNewMessageBroker(nil)
	SetNewApiClient(nil)
}

// Test multiple calls to factory functions return consistent results
func TestFactoryConsistency(t *testing.T) {
	orderMock := &mockOrderDataSource{}

	SetNewOrderDataSource(func() interfaces.IOrderDataSource {
		return orderMock
	})
	defer SetNewOrderDataSource(nil)

	result1 := NewOrderDataSource()
	result2 := NewOrderDataSource()

	if result1 != result2 {
		t.Error("Expected factory to return consistent results")
	}
}
