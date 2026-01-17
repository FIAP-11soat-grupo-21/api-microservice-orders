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

type OrderErrorMessage struct {
	OrderID         string `json:"order_id"`
	SystemTriggered string `json:"system_triggered"`
}

type MessageBroker interface {
	PublishOnTopic(ctx context.Context, topic string, message interface{}) error
	ConsumeOrderUpdates(ctx context.Context, handler OrderUpdateHandler) error
	ConsumeOrderError(ctx context.Context, handler OrderErrorHandler) error
	Close() error
}

type OrderUpdateHandler func(message OrderUpdateMessage) error

type OrderErrorHandler func(message OrderErrorMessage) error

type BrokerConfig struct {
	Type string

	// AWS
	AWSRegion          string
	AWSEndpoint        string
	AWSAccessKey       string
	AWSSecretAccessKey string

	// SQS
	SQSUpdateOrderStatusQueueURL string
	SQSOrderErrorQueueURL        string

	// SNS
	SNSOrderErrorTopicARN   string
	SNSOrderCreatedTopicARN string
}

type BrokerFactory interface {
	CreateBroker(config BrokerConfig) (MessageBroker, error)
}
