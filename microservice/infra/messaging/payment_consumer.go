package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"microservice/utils/config"
)

type PaymentEventMessage struct {
	OrderID       string  `json:"order_id"`
	PaymentStatus string  `json:"payment_status"`
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	TransactionID string  `json:"transaction_id"`
}

type PaymentEventConsumer struct {
	queueName string
}

func NewPaymentEventConsumer() *PaymentEventConsumer {
	cfg := config.LoadConfig()
	return &PaymentEventConsumer{
		queueName: cfg.RabbitMQ.PaymentQueue,
	}
}

func (c *PaymentEventConsumer) Start() error {
	channel := GetChannel()

	_, err := DeclareQueue(c.queueName)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := channel.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go c.processMessages(msgs)

	log.Printf("Payment event consumer started, listening on queue: %s", c.queueName)
	return nil
}

func (c *PaymentEventConsumer) processMessages(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		c.handleMessage(msg)
	}
}

func (c *PaymentEventConsumer) handleMessage(msg amqp.Delivery) {
	var paymentEvent PaymentEventMessage

	err := json.Unmarshal(msg.Body, &paymentEvent)
	if err != nil {
		log.Printf("Error unmarshaling payment event: %v", err)
		msg.Nack(false, false)
		return
	}

	fmt.Println("========================================")
	fmt.Println("PAYMENT EVENT RECEIVED VIA RABBITMQ")
	fmt.Println("========================================")
	fmt.Printf("Order ID: %s\n", paymentEvent.OrderID)
	fmt.Printf("Payment Status: %s\n", paymentEvent.PaymentStatus)
	fmt.Printf("Payment Method: %s\n", paymentEvent.PaymentMethod)
	fmt.Printf("Amount: %.2f\n", paymentEvent.Amount)
	fmt.Printf("Transaction ID: %s\n", paymentEvent.TransactionID)
	fmt.Println("========================================")

	msg.Ack(false)
	log.Printf("Payment event processed successfully for order: %s", paymentEvent.OrderID)
}
