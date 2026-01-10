package use_cases

import (
	"microservice/internal/domain/entities"
	"microservice/internal/interfaces"
)

type FindAllOrderStatusUseCase struct {
	orderStatusGateway interfaces.IOrderStatusGateway
}

func NewFindAllOrderStatusUseCase(orderStatusGateway interfaces.IOrderStatusGateway) *FindAllOrderStatusUseCase {
	return &FindAllOrderStatusUseCase{
		orderStatusGateway: orderStatusGateway,
	}
}

func (uc *FindAllOrderStatusUseCase) Execute() ([]entities.OrderStatus, error) {
	return uc.orderStatusGateway.FindAll()
}
