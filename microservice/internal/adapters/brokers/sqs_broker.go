package brokers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSBroker struct {
	client         *sqs.Client
	ordersQueueURL string
}

func NewSQSBroker(brokerConfig BrokerConfig) (*SQSBroker, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(brokerConfig.AWSRegion))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	if brokerConfig.SQSOrdersQueueURL == "" {
		return nil, fmt.Errorf("SQS orders queue URL is required")
	}

	log.Printf("SQS: Configured with orders queue: %s", brokerConfig.SQSOrdersQueueURL)

	return &SQSBroker{
		client:         sqs.NewFromConfig(cfg),
		ordersQueueURL: brokerConfig.SQSOrdersQueueURL,
	}, nil
}

func (s *SQSBroker) ConsumeOrderUpdates(ctx context.Context, handler OrderUpdateHandler) error {
	log.Printf("SQS: Starting order updates consumer on orders queue: %s", s.ordersQueueURL)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("SQS: Stopping order updates consumer")
				return
			default:
				if err := s.pollOrderUpdateMessages(ctx, handler); err != nil {
					log.Printf("SQS: Error polling order update messages: %v", err)
					time.Sleep(5 * time.Second)
				}
			}
		}
	}()

	return nil
}

func (s *SQSBroker) pollOrderUpdateMessages(ctx context.Context, handler OrderUpdateHandler) error {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:              &s.ordersQueueURL,
		MaxNumberOfMessages:   10,
		WaitTimeSeconds:       10, // Long polling
		MessageAttributeNames: []string{"All"},
	}

	result, err := s.client.ReceiveMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to receive order update messages: %w", err)
	}

	for _, message := range result.Messages {
		if err := s.processOrderUpdateMessage(ctx, message, handler); err != nil {
			log.Printf("SQS: Error processing order update message: %v", err)
			continue
		}

		if err := s.deleteOrderUpdateMessage(ctx, message); err != nil {
			log.Printf("SQS: Error deleting order update message: %v", err)
		}
	}

	return nil
}

func (s *SQSBroker) processOrderUpdateMessage(ctx context.Context, message types.Message, handler OrderUpdateHandler) error {
	var updateMsg OrderUpdateMessage
	if err := json.Unmarshal([]byte(*message.Body), &updateMsg); err != nil {
		return fmt.Errorf("failed to unmarshal order update message: %w", err)
	}

	log.Printf("SQS: Processing order update for order %s", updateMsg.OrderID)

	return handler(updateMsg)
}

func (s *SQSBroker) deleteOrderUpdateMessage(ctx context.Context, message types.Message) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      &s.ordersQueueURL,
		ReceiptHandle: message.ReceiptHandle,
	}

	_, err := s.client.DeleteMessage(ctx, input)
	return err
}

func (s *SQSBroker) Close() error {
	return nil
}
