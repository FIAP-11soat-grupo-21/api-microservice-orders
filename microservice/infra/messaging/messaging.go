package messaging

import (
	"log"
	"strings"

	"microservice/internal/adapters/brokers"
	"microservice/utils/config"
)

var (
	broker brokers.MessageBroker
)

func Connect() error {
	cfg := config.LoadConfig()

	log.Printf("Connecting to message broker: %s", cfg.MessageBroker.Type)

	brokerConfig := brokers.BrokerConfig{
		Type: cfg.MessageBroker.Type,
		// SQS
		SQSOrdersQueueURL: cfg.MessageBroker.SQS.OrdersQueueURL,
		AWSRegion:         cfg.MessageBroker.SQS.AWSRegion,
		// RabbitMQ
		RabbitMQURL:         buildRabbitMQURL(cfg),
		RabbitMQOrdersQueue: cfg.MessageBroker.RabbitMQ.OrdersQueue,
	}

	if cfg.MessageBroker.Type == "sqs" {
		log.Printf("SQS Config - Orders Queue: %s", brokerConfig.SQSOrdersQueueURL)
		log.Printf("SQS Config - AWS Region: %s", brokerConfig.AWSRegion)
	}

	factory := brokers.NewFactory()
	var err error
	broker, err = factory.CreateBroker(brokerConfig)
	if err != nil {
		log.Printf("Failed to create message broker: %v", err)
		return err
	}

	log.Printf("Connected to %s successfully", strings.ToUpper(cfg.MessageBroker.Type))
	return nil
}

func GetBroker() brokers.MessageBroker {
	return broker
}

func buildRabbitMQURL(cfg *config.Config) string {
	return cfg.MessageBroker.RabbitMQ.URL
}

func Close() {
	if broker != nil {
		if err := broker.Close(); err != nil {
			log.Printf("Error closing broker: %v", err)
		}
	}
	log.Println("Message broker connection closed")
}
