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
	conn        *amqp.Connection
	channel     *amqp.Channel
	ordersQueue string
}

func NewRabbitMQBroker(brokerConfig BrokerConfig) (*RabbitMQBroker, error) {
	// Validate configuration
	if brokerConfig.RabbitMQURL == "" {
		return nil, fmt.Errorf("RabbitMQ URL is required")
	}
	if brokerConfig.RabbitMQOrdersQueue == "" {
		return nil, fmt.Errorf("RabbitMQ orders queue is required")
	}

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
		conn:        conn,
		channel:     channel,
		ordersQueue: brokerConfig.RabbitMQOrdersQueue,
	}

	if err := broker.declareQueues(); err != nil {
		broker.Close()
		return nil, err
	}

	return broker, nil
}

func (r *RabbitMQBroker) declareQueues() error {
	// Declarar fila de atualizações de pedidos
	_, err := r.channel.QueueDeclare(
		r.ordersQueue, // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare orders queue: %w", err)
	}

	return nil
}

func (r *RabbitMQBroker) ConsumeOrderUpdates(ctx context.Context, handler OrderUpdateHandler) error {
	msgs, err := r.channel.Consume(
		r.ordersQueue, // queue
		"",            // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return fmt.Errorf("failed to register order updates consumer: %w", err)
	}

	log.Printf("RabbitMQ: Consuming order updates from queue: %s", r.ordersQueue)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("RabbitMQ: Stopping order updates consumer")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("RabbitMQ: Order updates channel closed")
					return
				}

				if err := r.processOrderUpdateMessage(msg, handler); err != nil {
					log.Printf("RabbitMQ: Error processing order update message: %v", err)
					if nackErr := msg.Nack(false, false); nackErr != nil {
						log.Printf("RabbitMQ: Error nacking order update message: %v", nackErr)
					}
				} else {
					if ackErr := msg.Ack(false); ackErr != nil {
						log.Printf("RabbitMQ: Error acknowledging order update message: %v", ackErr)
					}
				}
			}
		}
	}()

	return nil
}

func (r *RabbitMQBroker) processOrderUpdateMessage(msg amqp.Delivery, handler OrderUpdateHandler) error {
	var updateMsg OrderUpdateMessage
	if err := json.Unmarshal(msg.Body, &updateMsg); err != nil {
		return fmt.Errorf("failed to unmarshal order update message: %w", err)
	}

	log.Printf("RabbitMQ: Processing order update for order %s", updateMsg.OrderID)

	return handler(updateMsg)
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

func (r *RabbitMQBroker) ConsumeOrderError(ctx context.Context, handler OrderErrorHandler) error {
	// not implemented yet
	return nil
}

func (r *RabbitMQBroker) PublishOnTopic(ctx context.Context, topic string, message interface{}) error {
	// not implemented yet
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
