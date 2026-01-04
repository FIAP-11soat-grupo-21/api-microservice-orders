package brokers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSBroker struct {
	client           *sqs.Client
	paymentQueueURL  string
	kitchenQueueURL  string
}

func NewSQSBroker(brokerConfig BrokerConfig) (*SQSBroker, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(brokerConfig.AWSRegion))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	if brokerConfig.SQSPaymentQueueURL == "" {
		return nil, fmt.Errorf("SQS payment queue URL is required")
	}
	if brokerConfig.SQSKitchenQueueURL == "" {
		return nil, fmt.Errorf("SQS kitchen queue URL is required")
	}

	log.Printf("SQS: Configured with payment queue: %s", brokerConfig.SQSPaymentQueueURL)
	log.Printf("SQS: Configured with kitchen queue: %s", brokerConfig.SQSKitchenQueueURL)

	return &SQSBroker{
		client:          sqs.NewFromConfig(cfg),
		paymentQueueURL: brokerConfig.SQSPaymentQueueURL,
		kitchenQueueURL: brokerConfig.SQSKitchenQueueURL,
	}, nil
}

func (s *SQSBroker) ConsumePaymentConfirmations(ctx context.Context, handler PaymentConfirmationHandler) error {
	log.Printf("SQS: Starting payment confirmation consumer on payment queue: %s", s.paymentQueueURL)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("SQS: Stopping payment confirmation consumer")
				return
			default:
				if err := s.pollMessages(ctx, handler); err != nil {
					log.Printf("SQS: Error polling messages: %v", err)
					time.Sleep(5 * time.Second)
				}
			}
		}
	}()

	return nil
}

func (s *SQSBroker) pollMessages(ctx context.Context, handler PaymentConfirmationHandler) error {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:              &s.paymentQueueURL, // Usar fila de pagamento
		MaxNumberOfMessages:   10,
		WaitTimeSeconds:       10, // Long polling
		MessageAttributeNames: []string{"All"},
	}

	result, err := s.client.ReceiveMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to receive messages: %w", err)
	}

	for _, message := range result.Messages {
		if err := s.processMessage(ctx, message, handler); err != nil {
			log.Printf("SQS: Error processing message: %v", err)
			continue
		}

		if err := s.deleteMessage(ctx, message); err != nil {
			log.Printf("SQS: Error deleting message: %v", err)
		}
	}

	return nil
}

func (s *SQSBroker) processMessage(ctx context.Context, message types.Message, handler PaymentConfirmationHandler) error {
	messageType := ""
	if attr, exists := message.MessageAttributes["message_type"]; exists && attr.StringValue != nil {
		messageType = *attr.StringValue
	}

	if messageType != "payment.confirmed" && messageType != "payment.failed" {
		log.Printf("SQS: Skipping non-payment message type: %s", messageType)
		return nil
	}

	var paymentMsg PaymentConfirmationMessage
	if err := json.Unmarshal([]byte(*message.Body), &paymentMsg); err != nil {
		return fmt.Errorf("failed to unmarshal payment message: %w", err)
	}

	log.Printf("SQS: Processing payment confirmation for order %s", paymentMsg.OrderID)

	return handler(paymentMsg)
}

func (s *SQSBroker) deleteMessage(ctx context.Context, message types.Message) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      &s.paymentQueueURL, // Usar fila de pagamento
		ReceiptHandle: message.ReceiptHandle,
	}

	_, err := s.client.DeleteMessage(ctx, input)
	return err
}

func (s *SQSBroker) SendToKitchen(message map[string]interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal kitchen message: %w", err)
	}

	input := &sqs.SendMessageInput{
		QueueUrl:    &s.kitchenQueueURL, // Usar fila da cozinha
		MessageBody: aws.String(string(body)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"message_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String(fmt.Sprintf("%v", message["type"])),
			},
			"order_id": {
				DataType:    aws.String("String"),
				StringValue: aws.String(fmt.Sprintf("%v", message["order_id"])),
			},
			"source": {
				DataType:    aws.String("String"),
				StringValue: aws.String("orders-api"),
			},
		},
	}

	result, err := s.client.SendMessage(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to send kitchen message: %w", err)
	}

	log.Printf("SQS: Sent message to kitchen queue %s for order %v (MessageId: %s)",
		s.kitchenQueueURL, message["order_id"], *result.MessageId)
	return nil
}

func (s *SQSBroker) Close() error {
	return nil
}
