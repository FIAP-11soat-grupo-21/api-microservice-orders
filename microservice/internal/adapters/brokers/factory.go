package brokers

import (
	"fmt"
	"strings"
)

type Factory struct{}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) CreateBroker(config BrokerConfig) (MessageBroker, error) {
	switch strings.ToLower(config.Type) {
	case "sqs":
		return NewSQSBroker(config)
	default:
		return nil, fmt.Errorf("unsupported broker type: %s", config.Type)
	}
}
