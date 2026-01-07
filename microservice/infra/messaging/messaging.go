package messaging

import (
	"context"
	"fmt"
	"log"
	"strings"

	"microservice/internal/adapters/brokers"
	"microservice/internal/adapters/consumers"
	"microservice/internal/use_cases"
	"microservice/utils/config"
)

var (
	paymentConsumer *consumers.PaymentConsumer
	broker          brokers.MessageBroker
)

func Connect() error {
	cfg := config.LoadConfig()

	log.Printf("Connecting to message broker: %s", cfg.MessageBroker.Type)

	brokerConfig := brokers.BrokerConfig{
		Type: cfg.MessageBroker.Type,
		// SQS - Duas filas separadas
		SQSPaymentQueueURL: cfg.MessageBroker.SQS.PaymentQueueURL,
		SQSKitchenQueueURL: cfg.MessageBroker.SQS.KitchenQueueURL,
		AWSRegion:          cfg.MessageBroker.SQS.AWSRegion,
		// RabbitMQ
		RabbitMQURL:          buildRabbitMQURL(cfg),
		RabbitMQPaymentQueue: cfg.MessageBroker.RabbitMQ.PaymentQueue,
		RabbitMQKitchenQueue: cfg.MessageBroker.RabbitMQ.KitchenQueue,
	}

	factory := brokers.NewFactory()
	var err error
	broker, err = factory.CreateBroker(brokerConfig)
	if err != nil {
		return err
	}

	log.Printf("Connected to %s successfully", strings.ToUpper(cfg.MessageBroker.Type))
	return nil
}

func GetBroker() brokers.MessageBroker {
	return broker
}

func SetupPaymentConsumer(processPaymentUC *use_cases.ProcessPaymentConfirmationUseCase) error {
	if broker == nil {
		return fmt.Errorf("broker not connected, call Connect() first")
	}

	paymentConsumer = consumers.NewPaymentConsumer(
		broker,
		processPaymentUC,
		broker,
	)

	ctx := context.Background()
	go func() {
		if err := paymentConsumer.Start(ctx); err != nil {
			log.Printf("Error starting payment consumer: %v", err)
		}
	}()

	log.Println("Payment consumer setup completed")
	return nil
}

func NewPaymentEventConsumer() *PaymentEventConsumer {
	return &PaymentEventConsumer{}
}

type PaymentEventConsumer struct{}

func (c *PaymentEventConsumer) Start() error {
	log.Println("PaymentEventConsumer.Start() called - consumer will be setup when dependencies are available")
	return nil
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
