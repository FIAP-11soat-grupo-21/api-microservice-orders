package use_cases

import (
	"microservice/internal/adapters/gateways"
	"microservice/internal/domain/entities"
)

type FindAllOrderStatusUseCase struct {
	orderStatusGateway gateways.OrderStatusGateway
}

func NewFindAllOrderStatusUseCase(orderStatusGateway gateways.OrderStatusGateway) *FindAllOrderStatusUseCase {
	return &FindAllOrderStatusUseCase{
		orderStatusGateway: orderStatusGateway,
	}
}

func (uc *FindAllOrderStatusUseCase) Execute() ([]entities.OrderStatus, error) {
	return uc.orderStatusGateway.FindAll()
}
