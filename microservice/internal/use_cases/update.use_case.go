package use_cases

import (
	"time"

	"microservice/internal/adapters/dtos"
	"microservice/internal/adapters/gateways"
	"microservice/internal/domain/entities"
	"microservice/internal/domain/exceptions"
)

type UpdateOrderUseCase struct {
	orderGateway       gateways.OrderGateway
	orderStatusGateway gateways.OrderStatusGateway
}

func NewUpdateOrderUseCase(orderGateway gateways.OrderGateway, orderStatusGateway gateways.OrderStatusGateway) *UpdateOrderUseCase {
	return &UpdateOrderUseCase{
		orderGateway:       orderGateway,
		orderStatusGateway: orderStatusGateway,
	}
}

func (uc *UpdateOrderUseCase) Execute(dto dtos.UpdateOrderDTO) (entities.Order, error) {
	err := entities.ValidateID(dto.ID)
	if err != nil {
		return entities.Order{}, err
	}

	order, err := uc.orderGateway.FindByID(dto.ID)
	if err != nil {
		return entities.Order{}, &exceptions.OrderNotFoundException{}
	}

	status, err := uc.orderStatusGateway.FindByID(dto.StatusID)
	if err != nil {
		return entities.Order{}, &exceptions.OrderStatusNotFoundException{}
	}

	order.Status = *status
	now := time.Now()
	order.UpdatedAt = &now

	err = uc.orderGateway.Update(*order)
	if err != nil {
		return entities.Order{}, err
	}

	return *order, nil
}
