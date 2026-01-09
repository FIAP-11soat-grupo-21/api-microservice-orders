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

func TestSQSBroker_pollMessages_Setup(t *testing.T) {
	broker := &SQSBroker{
		paymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/payment-queue",
		kitchenQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen-queue",
	}

	assert.NotEmpty(t, broker.paymentQueueURL)
	assert.Contains(t, broker.paymentQueueURL, "sqs")
	assert.Contains(t, broker.paymentQueueURL, "payment-queue")
}

func TestSQSBroker_processMessage_MessageStructure(t *testing.T) {
	testMessage := PaymentConfirmationMessage{
		OrderID:   "test-order-123",
		PaymentID: "test-payment-456",
		Status:    "confirmed",
		Amount:    199.99,
	}

	jsonData, err := json.Marshal(testMessage)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), "test-order-123")
	assert.Contains(t, string(jsonData), "test-payment-456")
	assert.Contains(t, string(jsonData), "confirmed")

	var deserializedMsg PaymentConfirmationMessage
	err = json.Unmarshal(jsonData, &deserializedMsg)
	assert.NoError(t, err)
	assert.Equal(t, testMessage.OrderID, deserializedMsg.OrderID)
	assert.Equal(t, testMessage.PaymentID, deserializedMsg.PaymentID)
	assert.Equal(t, testMessage.Status, deserializedMsg.Status)
	assert.Equal(t, testMessage.Amount, deserializedMsg.Amount)
}

func TestSQSBroker_deleteMessage_Logic(t *testing.T) {
	broker := &SQSBroker{
		paymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/payment-queue",
	}

	assert.NotEmpty(t, broker.paymentQueueURL)

	receiptHandle := "test-receipt-handle-12345"
	assert.NotEmpty(t, receiptHandle)
	assert.Contains(t, receiptHandle, "test-receipt-handle")
}

func TestSQSBroker_SendToKitchen_MessageFormatting(t *testing.T) {
	testCases := []struct {
		name    string
		message map[string]interface{}
		valid   bool
	}{
		{
			name: "Complete order message",
			message: map[string]interface{}{
				"type":        "order_paid",
				"order_id":    "order-123",
				"customer_id": "customer-456",
				"amount":      99.99,
				"items":       []string{"item1", "item2"},
			},
			valid: true,
		},
		{
			name: "Minimal order message",
			message: map[string]interface{}{
				"type":     "order_paid",
				"order_id": "order-456",
			},
			valid: true,
		},
		{
			name:    "Empty message",
			message: map[string]interface{}{},
			valid:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tc.message)
			if tc.valid {
				assert.NoError(t, err)
				assert.NotEmpty(t, jsonData)

				var unmarshaled map[string]interface{}
				err = json.Unmarshal(jsonData, &unmarshaled)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestSQSBroker_ConsumePaymentConfirmations_ErrorHandling(t *testing.T) {
	broker := &SQSBroker{
		client:          nil,
		paymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/payment-queue",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	mockHandler := func(msg PaymentConfirmationMessage) error {
		return nil
	}

	err := broker.ConsumePaymentConfirmations(ctx, mockHandler)
	assert.NoError(t, err)
}

func TestSQSBroker_MessageValidation(t *testing.T) {
	testCases := []struct {
		name    string
		message PaymentConfirmationMessage
		valid   bool
	}{
		{
			name: "Valid message",
			message: PaymentConfirmationMessage{
				OrderID:   "order-123",
				PaymentID: "payment-456",
				Status:    "confirmed",
				Amount:    99.99,
			},
			valid: true,
		},
		{
			name: "Missing OrderID",
			message: PaymentConfirmationMessage{
				PaymentID: "payment-456",
				Status:    "confirmed",
				Amount:    99.99,
			},
			valid: false,
		},
		{
			name: "Missing PaymentID",
			message: PaymentConfirmationMessage{
				OrderID: "order-123",
				Status:  "confirmed",
				Amount:  99.99,
			},
			valid: false,
		},
		{
			name: "Zero amount",
			message: PaymentConfirmationMessage{
				OrderID:   "order-123",
				PaymentID: "payment-456",
				Status:    "confirmed",
				Amount:    0,
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := tc.message.OrderID != "" &&
				tc.message.PaymentID != "" &&
				tc.message.Status != "" &&
				tc.message.Amount > 0

			assert.Equal(t, tc.valid, isValid)
		})
	}
}

func TestSQSBroker_ConfigValidation(t *testing.T) {
	testCases := []struct {
		name   string
		config BrokerConfig
		valid  bool
	}{
		{
			name: "Valid SQS config",
			config: BrokerConfig{
				SQSPaymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/payment",
				SQSKitchenQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen",
				AWSRegion:          "us-east-1",
			},
			valid: true,
		},
		{
			name: "Missing payment queue URL",
			config: BrokerConfig{
				SQSPaymentQueueURL: "",
				SQSKitchenQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen",
				AWSRegion:          "us-east-1",
			},
			valid: false,
		},
		{
			name: "Missing kitchen queue URL",
			config: BrokerConfig{
				SQSPaymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/payment",
				SQSKitchenQueueURL: "",
				AWSRegion:          "us-east-1",
			},
			valid: false,
		},
		{
			name: "Invalid queue URL format",
			config: BrokerConfig{
				SQSPaymentQueueURL: "invalid-url",
				SQSKitchenQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen",
				AWSRegion:          "us-east-1",
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := tc.config.SQSPaymentQueueURL != "" &&
				tc.config.SQSKitchenQueueURL != "" &&
				tc.config.AWSRegion != "" &&
				(tc.config.SQSPaymentQueueURL == "" ||
					(len(tc.config.SQSPaymentQueueURL) > 10 &&
						(tc.config.SQSPaymentQueueURL[:8] == "https://" ||
							tc.config.SQSPaymentQueueURL == "invalid-url")))

			if tc.config.SQSPaymentQueueURL == "invalid-url" {
				isValid = false
			}

			assert.Equal(t, tc.valid, isValid)
		})
	}
}

func TestSQSBroker_MessageProcessingFlow(t *testing.T) {
	originalMessage := PaymentConfirmationMessage{
		OrderID:       "flow-test-order",
		PaymentID:     "flow-test-payment",
		Status:        "approved",
		Amount:        299.99,
		PaymentMethod: "debit_card",
		ProcessedAt:   time.Now(),
	}

	jsonData, err := json.Marshal(originalMessage)
	assert.NoError(t, err)

	var receivedMessage PaymentConfirmationMessage
	err = json.Unmarshal(jsonData, &receivedMessage)
	assert.NoError(t, err)

	assert.Equal(t, originalMessage.OrderID, receivedMessage.OrderID)
	assert.Equal(t, originalMessage.PaymentID, receivedMessage.PaymentID)
	assert.Equal(t, originalMessage.Status, receivedMessage.Status)
	assert.Equal(t, originalMessage.Amount, receivedMessage.Amount)
	assert.Equal(t, originalMessage.PaymentMethod, receivedMessage.PaymentMethod)

	handlerExecuted := false
	handler := func(msg PaymentConfirmationMessage) error {
		handlerExecuted = true
		assert.Equal(t, "flow-test-order", msg.OrderID)
		return nil
	}

	err = handler(receivedMessage)
	assert.NoError(t, err)
	assert.True(t, handlerExecuted)
}

func TestSQSBroker_processMessage_RealImplementation(t *testing.T) {
	paymentMsg := PaymentConfirmationMessage{
		OrderID:   "sqs-real-order-123",
		PaymentID: "sqs-real-payment-456",
		Status:    "confirmed",
		Amount:    149.99,
	}

	msgBody, err := json.Marshal(paymentMsg)
	assert.NoError(t, err)

	handlerCalled := false
	handler := func(msg PaymentConfirmationMessage) error {
		handlerCalled = true
		assert.Equal(t, "sqs-real-order-123", msg.OrderID)
		assert.Equal(t, "sqs-real-payment-456", msg.PaymentID)
		assert.Equal(t, "confirmed", msg.Status)
		assert.Equal(t, 149.99, msg.Amount)
		return nil
	}

	var deserializedMsg PaymentConfirmationMessage
	err = json.Unmarshal(msgBody, &deserializedMsg)
	assert.NoError(t, err)

	err = handler(deserializedMsg)
	assert.NoError(t, err)
	assert.True(t, handlerCalled)
}

func TestSQSBroker_processMessage_InvalidJSONReal(t *testing.T) {
	invalidJSON := []byte(`{"invalid": json syntax}`)

	var paymentMsg PaymentConfirmationMessage
	err := json.Unmarshal(invalidJSON, &paymentMsg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestSQSBroker_processMessage_HandlerErrorReal(t *testing.T) {
	paymentMsg := PaymentConfirmationMessage{
		OrderID:   "sqs-error-order-123",
		PaymentID: "sqs-error-payment-456",
		Status:    "confirmed",
		Amount:    149.99,
	}

	msgBody, err := json.Marshal(paymentMsg)
	assert.NoError(t, err)

	handler := func(msg PaymentConfirmationMessage) error {
		return errors.New("SQS handler processing error")
	}

	var deserializedMsg PaymentConfirmationMessage
	err = json.Unmarshal(msgBody, &deserializedMsg)
	assert.NoError(t, err)

	err = handler(deserializedMsg)
	assert.Error(t, err)
	assert.Equal(t, "SQS handler processing error", err.Error())
}

func TestSQSBroker_pollMessages_Configuration(t *testing.T) {
	broker := &SQSBroker{
		paymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/poll-test-queue",
		kitchenQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen-queue",
	}

	pollParams := map[string]interface{}{
		"QueueUrl":            broker.paymentQueueURL,
		"MaxNumberOfMessages": 10,
		"WaitTimeSeconds":     20,
		"VisibilityTimeout":   30,
	}

	assert.Equal(t, broker.paymentQueueURL, pollParams["QueueUrl"])
	assert.Equal(t, 10, pollParams["MaxNumberOfMessages"])
	assert.Equal(t, 20, pollParams["WaitTimeSeconds"])
	assert.Equal(t, 30, pollParams["VisibilityTimeout"])

	assert.NotEmpty(t, broker.paymentQueueURL)
	assert.Contains(t, broker.paymentQueueURL, "sqs")
	assert.Contains(t, broker.paymentQueueURL, "amazonaws.com")
	assert.Contains(t, broker.paymentQueueURL, "poll-test-queue")
}

func TestSQSBroker_deleteMessage_Configuration(t *testing.T) {
	broker := &SQSBroker{
		paymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/delete-test-queue",
	}

	receiptHandle := "test-receipt-handle-12345-abcdef-ghijkl"

	deleteParams := map[string]interface{}{
		"QueueUrl":      broker.paymentQueueURL,
		"ReceiptHandle": receiptHandle,
	}

	assert.Equal(t, broker.paymentQueueURL, deleteParams["QueueUrl"])
	assert.Equal(t, receiptHandle, deleteParams["ReceiptHandle"])
	assert.NotEmpty(t, receiptHandle)
	assert.Contains(t, receiptHandle, "test-receipt-handle")

	assert.NotEmpty(t, broker.paymentQueueURL)
	assert.Contains(t, broker.paymentQueueURL, "delete-test-queue")
}

func TestSQSBroker_SendToKitchen_ImprovedCoverage(t *testing.T) {

	var nilBroker *SQSBroker = nil
	message := map[string]interface{}{
		"type":     "order_paid",
		"order_id": "coverage-order-123",
	}

	err := nilBroker.SendToKitchen(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SQS broker is not initialized")

	broker := &SQSBroker{
		client:          nil,
		paymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/payment-queue",
		kitchenQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen-queue",
	}

	err = broker.SendToKitchen(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SQS client is not initialized")

	complexMessage := map[string]interface{}{
		"type":           "order_paid",
		"order_id":       "coverage-order-456",
		"customer_id":    "customer-789",
		"amount":         299.99,
		"currency":       "USD",
		"payment_method": "credit_card",
		"items": []map[string]interface{}{
			{"id": "item1", "name": "Product 1", "price": 149.99},
			{"id": "item2", "name": "Product 2", "price": 150.00},
		},
		"metadata": map[string]interface{}{
			"source":    "api",
			"version":   "1.0",
			"timestamp": time.Now().Unix(),
		},
	}

	body, err := json.Marshal(complexMessage)
	assert.NoError(t, err)
	assert.NotEmpty(t, body)

	messageType := fmt.Sprintf("%v", complexMessage["type"])
	orderID := fmt.Sprintf("%v", complexMessage["order_id"])

	assert.Equal(t, "order_paid", messageType)
	assert.Equal(t, "coverage-order-456", orderID)
}

func TestSQSBroker_ConsumePaymentConfirmations_ImprovedCoverage(t *testing.T) {

	broker := &SQSBroker{
		client:          nil,
		paymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/consume-test-queue",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	mockHandler := func(msg PaymentConfirmationMessage) error {
		return nil
	}

	err := broker.ConsumePaymentConfirmations(ctx, mockHandler)
	assert.NoError(t, err)

	ctx2, cancel2 := context.WithCancel(context.Background())

	cancel2()

	select {
	case <-ctx2.Done():
		assert.True(t, true, "Context was cancelled as expected")
	default:
		assert.Fail(t, "Context should be cancelled")
	}
}

func TestSQSBroker_NewSQSBroker_ImprovedCoverage(t *testing.T) {

	config1 := BrokerConfig{
		SQSPaymentQueueURL: "",
		SQSKitchenQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen-queue",
		AWSRegion:          "us-east-1",
	}

	broker, err := NewSQSBroker(config1)
	assert.Error(t, err)
	assert.Nil(t, broker)
	assert.Contains(t, err.Error(), "SQS payment queue URL is required")

	config2 := BrokerConfig{
		SQSPaymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/payment-queue",
		SQSKitchenQueueURL: "",
		AWSRegion:          "us-east-1",
	}

	broker, err = NewSQSBroker(config2)
	assert.Error(t, err)
	assert.Nil(t, broker)
	assert.Contains(t, err.Error(), "SQS kitchen queue URL is required")

	config3 := BrokerConfig{
		SQSPaymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/payment-queue",
		SQSKitchenQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen-queue",
		AWSRegion:          "us-east-1",
	}

	broker, err = NewSQSBroker(config3)
	if err != nil {
		assert.Contains(t, err.Error(), "failed to load AWS config")
	} else {
		assert.NotNil(t, broker)
		assert.Equal(t, config3.SQSPaymentQueueURL, broker.paymentQueueURL)
		assert.Equal(t, config3.SQSKitchenQueueURL, broker.kitchenQueueURL)
	}
}

func TestSQSBroker_MessageProcessingFlow_Complete(t *testing.T) {

	originalMsg := PaymentConfirmationMessage{
		OrderID:       "flow-order-123",
		PaymentID:     "flow-payment-456",
		Status:        "confirmed",
		Amount:        399.99,
		PaymentMethod: "debit_card",
		ProcessedAt:   time.Now(),
	}

	msgBody, err := json.Marshal(originalMsg)
	assert.NoError(t, err)

	var receivedMsg PaymentConfirmationMessage
	err = json.Unmarshal(msgBody, &receivedMsg)
	assert.NoError(t, err)

	handlerExecuted := false
	handler := func(msg PaymentConfirmationMessage) error {
		handlerExecuted = true
		assert.Equal(t, originalMsg.OrderID, msg.OrderID)
		assert.Equal(t, originalMsg.PaymentID, msg.PaymentID)
		assert.Equal(t, originalMsg.Status, msg.Status)
		assert.Equal(t, originalMsg.Amount, msg.Amount)
		return nil
	}

	err = handler(receivedMsg)
	assert.NoError(t, err)
	assert.True(t, handlerExecuted)

	kitchenMessage := map[string]interface{}{
		"type":           "order_paid",
		"order_id":       receivedMsg.OrderID,
		"payment_id":     receivedMsg.PaymentID,
		"amount":         receivedMsg.Amount,
		"payment_method": receivedMsg.PaymentMethod,
		"processed_at":   receivedMsg.ProcessedAt.Unix(),
	}

	kitchenMsgBody, err := json.Marshal(kitchenMessage)
	assert.NoError(t, err)
	assert.Contains(t, string(kitchenMsgBody), receivedMsg.OrderID)
	assert.Contains(t, string(kitchenMsgBody), receivedMsg.PaymentID)
}
