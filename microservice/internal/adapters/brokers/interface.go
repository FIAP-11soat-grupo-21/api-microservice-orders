package brokers

import (
	"context"
	"time"
)

type OrderUpdateMessage struct {
	Type      string                 `json:"type"`
	OrderID   string                 `json:"order_id"`
	Status    string                 `json:"status"`
	UpdatedAt time.Time              `json:"updated_at"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type MessageBroker interface {
	ConsumeOrderUpdates(ctx context.Context, handler OrderUpdateHandler) error
	Close() error
}

type OrderUpdateHandler func(message OrderUpdateMessage) error

type BrokerConfig struct {
	Type string

	// SQS Config
	SQSOrdersQueueURL string
	AWSRegion         string

	// RabbitMQ
	RabbitMQURL         string
	RabbitMQOrdersQueue string
}

type BrokerFactory interface {
	CreateBroker(config BrokerConfig) (MessageBroker, error)
}
