package mocks

import (
	"context"
	"microservice/internal/adapters/brokers"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMessageBroker is a mock of MessageBroker interface.
type MockMessageBroker struct {
	ctrl     *gomock.Controller
	recorder *MockMessageBrokerMockRecorder
}

// MockMessageBrokerMockRecorder is the mock recorder for MockMessageBroker.
type MockMessageBrokerMockRecorder struct {
	mock *MockMessageBroker
}

// NewMockMessageBroker creates a new mock instance.
func NewMockMessageBroker(ctrl *gomock.Controller) *MockMessageBroker {
	mock := &MockMessageBroker{ctrl: ctrl}
	mock.recorder = &MockMessageBrokerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageBroker) EXPECT() *MockMessageBrokerMockRecorder {
	return m.recorder
}

// PublishOnTopic mocks base method.
func (m *MockMessageBroker) PublishOnTopic(ctx context.Context, topic string, message interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishOnTopic", ctx, topic, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishOnTopic indicates an expected call of PublishOnTopic.
func (mr *MockMessageBrokerMockRecorder) PublishOnTopic(ctx, topic, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishOnTopic", reflect.TypeOf((*MockMessageBroker)(nil).PublishOnTopic), ctx, topic, message)
}

// ConsumeOrderUpdates mocks base method.
func (m *MockMessageBroker) ConsumeOrderUpdates(ctx context.Context, handler brokers.OrderUpdateHandler) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConsumeOrderUpdates", ctx, handler)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConsumeOrderUpdates indicates an expected call of ConsumeOrderUpdates.
func (mr *MockMessageBrokerMockRecorder) ConsumeOrderUpdates(ctx, handler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConsumeOrderUpdates", reflect.TypeOf((*MockMessageBroker)(nil).ConsumeOrderUpdates), ctx, handler)
}

// ConsumeOrderError mocks base method.
func (m *MockMessageBroker) ConsumeOrderError(ctx context.Context, handler brokers.OrderErrorHandler) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConsumeOrderError", ctx, handler)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConsumeOrderError indicates an expected call of ConsumeOrderError.
func (mr *MockMessageBrokerMockRecorder) ConsumeOrderError(ctx, handler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConsumeOrderError", reflect.TypeOf((*MockMessageBroker)(nil).ConsumeOrderError), ctx, handler)
}

// Close mocks base method.
func (m *MockMessageBroker) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockMessageBrokerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockMessageBroker)(nil).Close))
}
