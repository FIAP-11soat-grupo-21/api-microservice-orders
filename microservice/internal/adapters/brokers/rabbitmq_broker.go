package brokers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/streadway/amqp"
)

type RabbitMQBroker struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	paymentQueue string
	kitchenQueue string
}

func NewRabbitMQBroker(brokerConfig BrokerConfig) (*RabbitMQBroker, error) {
	conn, err := amqp.Dial(brokerConfig.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	broker := &RabbitMQBroker{
		conn:         conn,
		channel:      channel,
		paymentQueue: brokerConfig.RabbitMQPaymentQueue,
		kitchenQueue: brokerConfig.RabbitMQKitchenQueue,
	}

	if err := broker.declareQueues(); err != nil {
		broker.Close()
		return nil, err
	}

	return broker, nil
}

func (r *RabbitMQBroker) declareQueues() error {
	_, err := r.channel.QueueDeclare(
		r.paymentQueue, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare payment queue: %w", err)
	}

	// Declarar fila da cozinha
	_, err = r.channel.QueueDeclare(
		r.kitchenQueue, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare kitchen queue: %w", err)
	}

	return nil
}

func (r *RabbitMQBroker) ConsumePaymentConfirmations(ctx context.Context, handler PaymentConfirmationHandler) error {
	msgs, err := r.channel.Consume(
		r.paymentQueue, // queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("RabbitMQ: Consuming payment confirmations from queue: %s", r.paymentQueue)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("RabbitMQ: Stopping payment confirmation consumer")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("RabbitMQ: Channel closed")
					return
				}

				if err := r.processPaymentMessage(msg, handler); err != nil {
					log.Printf("RabbitMQ: Error processing message: %v", err)

					if r.shouldDiscardMessage(err) {
						log.Printf("RabbitMQ: Discarding message due to non-recoverable error: %v", err)
						if ackErr := msg.Ack(false); ackErr != nil {
							log.Printf("RabbitMQ: Error acknowledging message: %v", ackErr)
						}
					} else {
						log.Printf("RabbitMQ: Rejecting message without requeue")
						if nackErr := msg.Nack(false, false); nackErr != nil {
							log.Printf("RabbitMQ: Error nacking message: %v", nackErr)
						}
					}
				} else {
					if ackErr := msg.Ack(false); ackErr != nil {
						log.Printf("RabbitMQ: Error acknowledging message: %v", ackErr)
					}
				}
			}
		}
	}()

	return nil
}

func (r *RabbitMQBroker) processPaymentMessage(msg amqp.Delivery, handler PaymentConfirmationHandler) error {
	var paymentMsg PaymentConfirmationMessage
	if err := json.Unmarshal(msg.Body, &paymentMsg); err != nil {
		return fmt.Errorf("failed to unmarshal payment message: %w", err)
	}

	log.Printf("RabbitMQ: Processing payment confirmation for order %s", paymentMsg.OrderID)

	return handler(paymentMsg)
}

func (r *RabbitMQBroker) shouldDiscardMessage(err error) bool {
	errorMsg := err.Error()

	nonRecoverableErrors := []string{
		"Order not found",
		"Invalid order ID",
		"Invalid payment confirmation",
		"order ID is required",
		"payment ID is required",
		"payment status is required",
	}

	for _, nonRecoverableError := range nonRecoverableErrors {
		if strings.Contains(errorMsg, nonRecoverableError) {
			return true
		}
	}

	return false
}

func (r *RabbitMQBroker) SendToKitchen(message map[string]interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal kitchen message: %w", err)
	}

	err = r.channel.Publish(
		"",             // exchange
		r.kitchenQueue, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Headers: amqp.Table{
				"message_type": fmt.Sprintf("%v", message["type"]),
				"order_id":     fmt.Sprintf("%v", message["order_id"]),
			},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish kitchen message: %w", err)
	}

	log.Printf("RabbitMQ: Sent message to kitchen queue for order %v", message["order_id"])
	return nil
}

func (r *RabbitMQBroker) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	return nil
}
