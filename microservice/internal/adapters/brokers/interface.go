package brokers

import (
	"context"
	"time"
)

type PaymentConfirmationMessage struct {
	Type          string                 `json:"type"`
	OrderID       string                 `json:"order_id"`
	PaymentID     string                 `json:"payment_id"`
	Amount        float64                `json:"amount"`
	Status        string                 `json:"status"`
	ProcessedAt   time.Time              `json:"processed_at"`
	PaymentMethod string                 `json:"payment_method,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

type MessageBroker interface {
	ConsumePaymentConfirmations(ctx context.Context, handler PaymentConfirmationHandler) error

	SendToKitchen(message map[string]interface{}) error

	Close() error
}

type PaymentConfirmationHandler func(message PaymentConfirmationMessage) error

type BrokerConfig struct {
	Type string

	// SQS Config
	SQSPaymentQueueURL string
	SQSKitchenQueueURL string
	AWSRegion          string

	// RabbitMQ
	RabbitMQURL          string
	RabbitMQPaymentQueue string
	RabbitMQKitchenQueue string
}

type BrokerFactory interface {
	CreateBroker(config BrokerConfig) (MessageBroker, error)
}
