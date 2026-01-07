package use_cases

import (
	"log"
	"time"

	"microservice/internal/adapters/brokers"
	"microservice/internal/adapters/dtos"
	"microservice/internal/adapters/gateways"
	"microservice/internal/domain/entities"
	"microservice/internal/domain/exceptions"
	identityUtils "microservice/utils/identity"
)

const INITIAL_ORDER_STATUS_ID = "56d3b3c3-1801-49cd-bae7-972c78082012"

type CreateOrderUseCase struct {
	orderGateway       gateways.OrderGateway
	orderStatusGateway gateways.OrderStatusGateway
	messageBroker      brokers.MessageBroker
}

func NewCreateOrderUseCase(orderGateway gateways.OrderGateway, orderStatusGateway gateways.OrderStatusGateway, messageBroker brokers.MessageBroker) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderGateway:       orderGateway,
		orderStatusGateway: orderStatusGateway,
		messageBroker:      messageBroker,
	}
}

func (uc *CreateOrderUseCase) Execute(customerID *string, items []dtos.CreateOrderItemDTO) (entities.Order, error) {
	status, err := uc.orderStatusGateway.FindByID(INITIAL_ORDER_STATUS_ID)
	if err != nil {
		return entities.Order{}, &exceptions.OrderStatusNotFoundException{}
	}

	order, _ := entities.NewOrder(identityUtils.NewUUIDV4(), customerID)
	order.Status = *status
	order.CreatedAt = time.Now()

	for _, item := range items {
		orderItem, err := entities.NewOrderItem(
			identityUtils.NewUUIDV4(),
			item.ProductID,
			order.ID,
			item.Quantity,
			item.Price,
		)
		if err != nil {
			return entities.Order{}, err
		}
		order.AddItem(*orderItem)
	}

	err = order.CalcTotalAmount()
	if err != nil {
		return entities.Order{}, err
	}

	err = uc.orderGateway.Create(*order)
	if err != nil {
		return entities.Order{}, err
	}

	if uc.messageBroker != nil {
		items := make([]map[string]interface{}, len(order.Items))
		for i, item := range order.Items {
			items[i] = map[string]interface{}{
				"id":         item.ID,
				"product_id": item.ProductID.Value(),
				"quantity":   item.Quantity.Value(),
				"unit_price": item.UnitPrice.Value(),
			}
		}

		kitchenMessage := map[string]interface{}{
			"order_id":    order.ID,
			"customer_id": order.CustomerID,
			"amount":      order.Amount.Value(),
			"items":       items,
		}

		if err := uc.messageBroker.SendToKitchen(kitchenMessage); err != nil {
			log.Printf("Warning: Failed to send message to kitchen: %v", err)
		} else {
			log.Printf("Kitchen order creation message sent for order %s (Amount: %.2f, Items: %d)", order.ID, order.Amount.Value(), len(order.Items))
		}
	}

	return *order, nil
}
