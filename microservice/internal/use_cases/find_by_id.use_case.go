package use_cases

import (
	"microservice/internal/adapters/gateways"
	"microservice/internal/domain/entities"
	"microservice/internal/domain/exceptions"
)

type FindOrderByIDUseCase struct {
	orderGateway gateways.OrderGateway
}

func NewFindOrderByIDUseCase(orderGateway gateways.OrderGateway) *FindOrderByIDUseCase {
	return &FindOrderByIDUseCase{
		orderGateway: orderGateway,
	}
}

func (uc *FindOrderByIDUseCase) Execute(id string) (entities.Order, error) {
	err := entities.ValidateID(id)
	if err != nil {
		return entities.Order{}, err
	}

	order, err := uc.orderGateway.FindByID(id)
	if err != nil {
		return entities.Order{}, &exceptions.OrderNotFoundException{}
	}

	return *order, nil
}
