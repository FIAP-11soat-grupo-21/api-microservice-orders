package messaging

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"microservice/utils/config"
)

var (
	connection *amqp.Connection
	channel    *amqp.Channel
)

func GetConnection() *amqp.Connection {
	return connection
}

func GetChannel() *amqp.Channel {
	return channel
}

func Connect() error {
	cfg := config.LoadConfig()

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	var err error
	maxRetries := 10
	retryInterval := 3 * time.Second

	for i := range maxRetries {
		connection, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w", maxRetries, err)
	}

	channel, err = connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	log.Println("Connected to RabbitMQ successfully")
	return nil
}

func DeclareQueue(queueName string) (amqp.Queue, error) {
	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	return queue, nil
}

func Close() {
	if channel != nil {
		channel.Close()
	}
	if connection != nil {
		connection.Close()
	}
	log.Println("RabbitMQ connection closed")
}
