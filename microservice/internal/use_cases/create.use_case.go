package use_cases

import (
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

	return *order, nil
}
