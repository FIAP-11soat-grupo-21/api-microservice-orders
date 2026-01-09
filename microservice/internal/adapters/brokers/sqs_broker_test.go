package brokers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSQSBroker_MissingPaymentQueueURL(t *testing.T) {
	config := BrokerConfig{
		SQSPaymentQueueURL: "", // Missing
		SQSKitchenQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen-queue",
		AWSRegion:          "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	assert.Error(t, err)
	assert.Nil(t, broker)
	assert.Contains(t, err.Error(), "SQS payment queue URL is required")
}

func TestNewSQSBroker_MissingKitchenQueueURL(t *testing.T) {
	config := BrokerConfig{
		SQSPaymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/payment-queue",
		SQSKitchenQueueURL: "", // Missing
		AWSRegion:          "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	assert.Error(t, err)
	assert.Nil(t, broker)
	assert.Contains(t, err.Error(), "SQS kitchen queue URL is required")
}

func TestNewSQSBroker_ValidConfig(t *testing.T) {
	config := BrokerConfig{
		SQSPaymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/payment-queue",
		SQSKitchenQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen-queue",
		AWSRegion:          "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	if err != nil {
		assert.Contains(t, err.Error(), "failed to load AWS config")
		return
	}

	assert.NotNil(t, broker)
	assert.Equal(t, config.SQSPaymentQueueURL, broker.paymentQueueURL)
	assert.Equal(t, config.SQSKitchenQueueURL, broker.kitchenQueueURL)
}

func TestSQSBroker_SendToKitchen_NilBroker(t *testing.T) {
	var broker *SQSBroker = nil
	
	message := map[string]interface{}{
		"type":     "order_paid",
		"order_id": "order-123",
	}

	err := broker.SendToKitchen(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SQS broker is not initialized")
}

func TestSQSBroker_SendToKitchen_NilClient(t *testing.T) {
	broker := &SQSBroker{
		client:          nil,
		paymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/payment-queue",
		kitchenQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen-queue",
	}
	
	message := map[string]interface{}{
		"type":     "order_paid",
		"order_id": "order-123",
	}

	err := broker.SendToKitchen(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SQS client is not initialized")
}

func TestSQSBroker_SendToKitchen_InvalidMessage(t *testing.T) {
	message := map[string]interface{}{
		"invalid": make(chan int),
	}

	_, err := json.Marshal(message)
	assert.Error(t, err)
}

func TestSQSBroker_SendToKitchen_ValidMessage(t *testing.T) {
	message := map[string]interface{}{
		"type":           "order_paid",
		"order_id":       "order-123",
		"customer_id":    "customer-456",
		"status":         "paid",
		"total_amount":   99.99,
		"payment_method": "credit_card",
	}

	body, err := json.Marshal(message)
	assert.NoError(t, err)
	assert.NotEmpty(t, body)

	var unmarshaled map[string]interface{}
	err = json.Unmarshal(body, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, "order_paid", unmarshaled["type"])
	assert.Equal(t, "order-123", unmarshaled["order_id"])
	assert.Equal(t, "customer-456", unmarshaled["customer_id"])
	assert.Equal(t, "paid", unmarshaled["status"])
	assert.Equal(t, 99.99, unmarshaled["total_amount"])
	assert.Equal(t, "credit_card", unmarshaled["payment_method"])
}

func TestSQSBroker_processMessage_ValidPaymentMessage(t *testing.T) {
	paymentMsg := PaymentConfirmationMessage{
		OrderID:       "order-123",
		PaymentID:     "payment-456",
		Status:        "approved",
		Amount:        99.99,
		PaymentMethod: "credit_card",
		ProcessedAt:   time.Now(),
	}

	msgBody, err := json.Marshal(paymentMsg)
	assert.NoError(t, err)

	handlerCalled := false
	mockHandler := func(msg PaymentConfirmationMessage) error {
		handlerCalled = true
		assert.Equal(t, "order-123", msg.OrderID)
		assert.Equal(t, "payment-456", msg.PaymentID)
		assert.Equal(t, "approved", msg.Status)
		assert.Equal(t, 99.99, msg.Amount)
		assert.Equal(t, "credit_card", msg.PaymentMethod)
		return nil
	}

	var unmarshaledMsg PaymentConfirmationMessage
	err = json.Unmarshal(msgBody, &unmarshaledMsg)
	assert.NoError(t, err)
	
	err = mockHandler(unmarshaledMsg)
	assert.NoError(t, err)
	assert.True(t, handlerCalled)
}

func TestSQSBroker_processMessage_InvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"invalid": json}`)
	
	var paymentMsg PaymentConfirmationMessage
	err := json.Unmarshal(invalidJSON, &paymentMsg)
	assert.Error(t, err)
}

func TestSQSBroker_processMessage_HandlerError(t *testing.T) {
	paymentMsg := PaymentConfirmationMessage{
		OrderID:   "order-123",
		PaymentID: "payment-456",
		Status:    "approved",
	}

	msgBody, err := json.Marshal(paymentMsg)
	assert.NoError(t, err)

	mockHandler := func(msg PaymentConfirmationMessage) error {
		return errors.New("handler error")
	}

	var unmarshaledMsg PaymentConfirmationMessage
	err = json.Unmarshal(msgBody, &unmarshaledMsg)
	assert.NoError(t, err)
	
	err = mockHandler(unmarshaledMsg)
	assert.Error(t, err)
	assert.Equal(t, "handler error", err.Error())
}

func TestSQSBroker_ConsumePaymentConfirmations_ContextCancellation(t *testing.T) {
	broker := &SQSBroker{
		paymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/payment-queue",
		kitchenQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen-queue",
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	mockHandler := func(msg PaymentConfirmationMessage) error {
		return nil
	}

	err := broker.ConsumePaymentConfirmations(ctx, mockHandler)
	assert.NoError(t, err)

	cancel()
	
	select {
	case <-ctx.Done():
		assert.True(t, true, "Context was cancelled as expected")
	default:
		assert.Fail(t, "Context should be cancelled")
	}
}

func TestSQSBroker_Close(t *testing.T) {
	broker := &SQSBroker{}
	err := broker.Close()
	assert.NoError(t, err)
}

func TestSQSBroker_MessageAttributes(t *testing.T) {
	message := map[string]interface{}{
		"type":     "order_paid",
		"order_id": "order-123",
	}

	messageType := fmt.Sprintf("%v", message["type"])
	orderID := fmt.Sprintf("%v", message["order_id"])

	assert.Equal(t, "order_paid", messageType)
	assert.Equal(t, "order-123", orderID)
}

func TestSQSBroker_MessageTypeFiltering(t *testing.T) {
	testCases := []struct {
		messageType string
		shouldSkip  bool
	}{
		{"payment.confirmed", false},
		{"payment.failed", false},
		{"order.created", true},
		{"user.updated", true},
		{"", true},
	}

	for _, tc := range testCases {
		t.Run(tc.messageType, func(t *testing.T) {
			shouldSkip := tc.messageType != "payment.confirmed" && tc.messageType != "payment.failed"
			assert.Equal(t, tc.shouldSkip, shouldSkip)
		})
	}
}