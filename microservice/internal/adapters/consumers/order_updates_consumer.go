package consumers

import (
	"context"
	"log"

	"microservice/internal/adapters/brokers"
	"microservice/internal/use_cases"
	"microservice/internal/interfaces"
)

type OrderUpdatesConsumer struct {
	broker                   brokers.MessageBroker
	updateOrderStatusUseCase *use_cases.UpdateOrderStatusUseCase
}

func NewOrderUpdatesConsumer(broker brokers.MessageBroker, orderGateway interfaces.IOrderGateway, orderStatusGateway interfaces.IOrderStatusGateway) *OrderUpdatesConsumer {
	updateOrderStatusUseCase := use_cases.NewUpdateOrderStatusUseCase(orderGateway, orderStatusGateway)
	
	return &OrderUpdatesConsumer{
		broker:                   broker,
		updateOrderStatusUseCase: updateOrderStatusUseCase,
	}
}

func (c *OrderUpdatesConsumer) Start(ctx context.Context) error {
	log.Println("Starting order updates consumer...")

	return c.broker.ConsumeOrderUpdates(ctx, c.processOrderUpdate)
}

func (c *OrderUpdatesConsumer) processOrderUpdate(message brokers.OrderUpdateMessage) error {
	log.Printf("Processing order update for order %s: %s", message.OrderID, message.Status)

	// Criar DTO para o use case
	updateDTO := use_cases.UpdateOrderStatusDTO{
		OrderID: message.OrderID,
		Status:  message.Status,
	}

	// Executar a atualização do status
	result, err := c.updateOrderStatusUseCase.Execute(updateDTO)
	if err != nil {
		log.Printf("Error updating order %s status: %v", message.OrderID, err)
		return err
	}

	log.Printf("Order %s status successfully updated to: %s", result.Order.ID, result.Order.Status.Name)
	return nil
}