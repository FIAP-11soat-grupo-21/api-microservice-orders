package consumers

import (
	"context"
	"log"
	"microservice/internal/adapters/brokers"
	"microservice/internal/interfaces"
	"microservice/internal/use_cases"
)

type OrderErrorConsumer struct {
	broker             brokers.MessageBroker
	deleteOrderUseCase *use_cases.DeleteOrderUseCase
}

func NewOrderErrorConsumer(broker brokers.MessageBroker, orderGateway interfaces.IOrderGateway) *OrderErrorConsumer {
	deleteOrderUseCase := use_cases.NewDeleteOrderUseCase(orderGateway)

	return &OrderErrorConsumer{
		broker:             broker,
		deleteOrderUseCase: deleteOrderUseCase,
	}
}

func (c *OrderErrorConsumer) Start(ctx context.Context) error {
	log.Println("Starting order error consumer...")

	return c.broker.ConsumeOrderError(ctx, c.processOrderError)
}

func (c *OrderErrorConsumer) processOrderError(message brokers.OrderErrorMessage) error {
	log.Printf("Processing order error for order: %s", message.OrderID)

	err := c.deleteOrderUseCase.Execute(message.OrderID)

	if err != nil {
		return err
	}

	log.Printf("Order %s successfully deleted", message.OrderID)
	return nil
}
