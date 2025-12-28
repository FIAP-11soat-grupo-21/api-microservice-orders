package use_cases

import (
	"microservice/internal/adapters/gateways"
	"microservice/internal/domain/entities"
	"microservice/internal/domain/exceptions"
)

type DeleteOrderUseCase struct {
	orderGateway gateways.OrderGateway
}

func NewDeleteOrderUseCase(orderGateway gateways.OrderGateway) *DeleteOrderUseCase {
	return &DeleteOrderUseCase{
		orderGateway: orderGateway,
	}
}

func (uc *DeleteOrderUseCase) Execute(id string) error {
	err := entities.ValidateID(id)
	if err != nil {
		return err
	}

	_, err = uc.orderGateway.FindByID(id)
	if err != nil {
		return &exceptions.OrderNotFoundException{}
	}

	return uc.orderGateway.Delete(id)
}
