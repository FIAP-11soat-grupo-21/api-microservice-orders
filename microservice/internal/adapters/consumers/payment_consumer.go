package consumers

import (
	"context"
	"log"
	"time"

	"microservice/internal/adapters/brokers"
	"microservice/internal/domain/entities"
	"microservice/internal/use_cases"
)

type PaymentConsumer struct {
	broker           brokers.MessageBroker
	processPaymentUC *use_cases.ProcessPaymentConfirmationUseCase
	kitchenBroker    brokers.MessageBroker
}

func NewPaymentConsumer(
	broker brokers.MessageBroker,
	processPaymentUC *use_cases.ProcessPaymentConfirmationUseCase,
	kitchenBroker brokers.MessageBroker,
) *PaymentConsumer {
	return &PaymentConsumer{
		broker:           broker,
		processPaymentUC: processPaymentUC,
		kitchenBroker:    kitchenBroker,
	}
}

func (c *PaymentConsumer) Start(ctx context.Context) error {
	log.Println("Starting payment confirmation consumer...")

	return c.broker.ConsumePaymentConfirmations(ctx, c.processPaymentConfirmation)
}

func (c *PaymentConsumer) processPaymentConfirmation(message brokers.PaymentConfirmationMessage) error {
	log.Printf("Processing payment confirmation for order %s: %s", message.OrderID, message.Status)

	dto := use_cases.PaymentConfirmationDTO{
		OrderID:       message.OrderID,
		PaymentID:     message.PaymentID,
		Status:        message.Status,
		Amount:        message.Amount,
		PaymentMethod: message.PaymentMethod,
		ProcessedAt:   message.ProcessedAt,
	}

	result, err := c.processPaymentUC.Execute(dto)
	if err != nil {
		log.Printf("Error processing payment confirmation: %v", err)
		return err
	}

	log.Printf("Payment processing result: %s", result.Message)

	if result.ShouldNotifyKitchen {
		return c.sendToKitchen(&result.Order, message)
	}

	return nil
}

func (c *PaymentConsumer) sendToKitchen(order *entities.Order, message brokers.PaymentConfirmationMessage) error {
	kitchenMessage := map[string]interface{}{
		"type":           "order_paid",
		"order_id":       order.ID,
		"customer_id":    order.CustomerID,
		"items":          order.Items,
		"status":         "paid",
		"created_at":     order.CreatedAt,
		"paid_at":        time.Now(),
		"total_amount":   order.Amount.Value(),
		"payment_method": message.PaymentMethod,
	}

	if err := c.kitchenBroker.SendToKitchen(kitchenMessage); err != nil {
		log.Printf("Failed to send order to kitchen: %v", err)
		return nil
	}

	log.Printf("Order %s sent to kitchen for preparation", order.ID)
	return nil
}
