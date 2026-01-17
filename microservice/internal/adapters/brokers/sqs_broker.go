package brokers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSBroker struct {
	sqsClient                 *sqs.Client
	snsClient                 *sns.Client
	updateOrderStatusQueueURL string
	orderErrorQueueURL        string
}

type SNSNotification struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}

func NewSQSBroker(brokerConfig BrokerConfig) (*SQSBroker, error) {
	if brokerConfig.SQSUpdateOrderStatusQueueURL == "" {
		return nil, fmt.Errorf("SQS update order status queue URL is required")
	}

	if brokerConfig.SQSOrderErrorQueueURL == "" {
		return nil, fmt.Errorf("SQS order error queue URL is required")
	}

	log.Printf("[SQS] Configured with update order status queue: %s", brokerConfig.SQSUpdateOrderStatusQueueURL)
	log.Printf("[SQS] Configured with order error queue: %s", brokerConfig.SQSOrderErrorQueueURL)

	optFns := []func(*config.LoadOptions) error{
		config.WithRegion(brokerConfig.AWSRegion),
		config.WithBaseEndpoint(brokerConfig.AWSEndpoint),
	}

	if brokerConfig.AWSAccessKey != "" && brokerConfig.AWSSecretAccessKey != "" {
		optFns = append(optFns, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				brokerConfig.AWSAccessKey,
				brokerConfig.AWSSecretAccessKey,
				"",
			),
		))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), optFns...)

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &SQSBroker{
		sqsClient:                 sqs.NewFromConfig(cfg),
		snsClient:                 sns.NewFromConfig(cfg),
		updateOrderStatusQueueURL: brokerConfig.SQSUpdateOrderStatusQueueURL,
		orderErrorQueueURL:        brokerConfig.SQSOrderErrorQueueURL,
	}, nil
}

func (s *SQSBroker) ConsumeOrderUpdates(ctx context.Context, handler OrderUpdateHandler) error {
	log.Printf("[SQS] Starting order updates consumer: %s", s.updateOrderStatusQueueURL)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("[SQS] Stopping order updates consumer")
				return
			default:
				if err := s.pollOrderUpdateMessages(ctx, handler); err != nil {
					log.Printf("[SQS] Error polling order update messages: %v", err)
					time.Sleep(5 * time.Second)
				}
			}
		}
	}()

	return nil
}

func (s *SQSBroker) pollOrderUpdateMessages(ctx context.Context, handler OrderUpdateHandler) error {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(s.updateOrderStatusQueueURL),
		MaxNumberOfMessages:   10,
		WaitTimeSeconds:       10, // Long polling
		MessageAttributeNames: []string{"All"},
	}

	result, err := s.sqsClient.ReceiveMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to receive order update messages: %w", err)
	}

	for _, message := range result.Messages {
		var updateMsg OrderUpdateMessage
		if err := s.unmarshalMessage(message, &updateMsg); err != nil {
			log.Printf("[SQS] failed to unmarshall order update message: %v", err)
			continue
		}

		log.Printf("[SQS] Processing order update for order %s", updateMsg.OrderID)

		err := handler(updateMsg)
		if err != nil {
			log.Printf("[SQS] Error processing order update message for order %s: %v", updateMsg.OrderID, err)
			continue
		}

		if err := s.deleteMessage(ctx, s.updateOrderStatusQueueURL, message); err != nil {
			log.Printf("[SQS] Error deleting order update message: %v", err)
		}
	}

	return nil
}

func (s *SQSBroker) ConsumeOrderError(ctx context.Context, handler OrderErrorHandler) error {
	log.Printf("[SQS] Starting order error consumer: %s", s.orderErrorQueueURL)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("[SQS] Stopping order error consumer")
				return
			default:
				if err := s.pollOrderErrorMessages(ctx, handler); err != nil {
					log.Printf("[SQS] Error polling order error messages: %v", err)
					time.Sleep(5 * time.Second)
				}
			}
		}
	}()

	return nil
}

func (s *SQSBroker) pollOrderErrorMessages(ctx context.Context, handler OrderErrorHandler) error {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:              &s.orderErrorQueueURL,
		MaxNumberOfMessages:   10,
		WaitTimeSeconds:       10, // Long polling
		MessageAttributeNames: []string{"All"},
	}

	result, err := s.sqsClient.ReceiveMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to receive order error messages: %w", err)
	}

	for _, message := range result.Messages {
		var orderErrorMsg OrderErrorMessage
		if err := s.unmarshalMessage(message, &orderErrorMsg); err != nil {
			log.Printf("[SQS] failed to unmarshall order error message: %v", err)
			continue
		}

		log.Printf("[SQS] Processing order error for order %s", orderErrorMsg.OrderID)

		err := handler(orderErrorMsg)
		if err != nil {
			log.Printf("[SQS] Error processing order error message for order %s: %v", orderErrorMsg.OrderID, err)
			s.deleteMessage(ctx, s.orderErrorQueueURL, message)
			continue
		}

		if err := s.deleteMessage(ctx, s.orderErrorQueueURL, message); err != nil {
			log.Printf("[SQS] Error deleting order error message: %v", err)
		}
	}

	return nil
}

func (s *SQSBroker) unmarshalMessage(message types.Message, obj any) error {
	var snsNotification SNSNotification

	if err := json.Unmarshal([]byte(*message.Body), &snsNotification); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(snsNotification.Message), &obj); err != nil {
		return err
	}

	return nil
}

func (s *SQSBroker) deleteMessage(ctx context.Context, queueURL string, message types.Message) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: message.ReceiptHandle,
	}

	_, err := s.sqsClient.DeleteMessage(ctx, input)
	return err
}

func (s *SQSBroker) Close() error {
	return nil
}

func (s *SQSBroker) PublishOnTopic(ctx context.Context, topic string, message interface{}) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	log.Printf("[SQS] Publishing message to topic: %s", topic)
	log.Printf("[SQS] Message payload: %s", string(payload))

	input := &sns.PublishInput{
		TopicArn: aws.String(topic),
		Message:  aws.String(string(payload)),
	}

	_, err = s.snsClient.Publish(ctx, input)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
