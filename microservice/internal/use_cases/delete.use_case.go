package use_cases

import (
	"microservice/internal/domain/entities"
	"microservice/internal/domain/exceptions"
	"microservice/internal/interfaces"
)

type DeleteOrderUseCase struct {
	orderGateway interfaces.IOrderGateway
}

func NewDeleteOrderUseCase(orderGateway interfaces.IOrderGateway) *DeleteOrderUseCase {
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
