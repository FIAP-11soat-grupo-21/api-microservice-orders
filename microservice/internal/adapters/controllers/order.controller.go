package controllers

import (
	"microservice/internal/adapters/brokers"
	"microservice/internal/adapters/dtos"
	"microservice/internal/adapters/gateways"
	"microservice/internal/adapters/presenters"
	"microservice/internal/interfaces"
	"microservice/internal/use_cases"
)

type OrderController struct {
	orderDataSource       interfaces.IOrderDataSource
	orderStatusDataSource interfaces.IOrderStatusDataSource
	orderGateway          gateways.OrderGateway
	orderStatusGateway    gateways.OrderStatusGateway
	messageBroker         brokers.MessageBroker
}

func NewOrderController(orderDataSource interfaces.IOrderDataSource, orderStatusDataSource interfaces.IOrderStatusDataSource, messageBroker brokers.MessageBroker) *OrderController {
	return &OrderController{
		orderDataSource:       orderDataSource,
		orderStatusDataSource: orderStatusDataSource,
		orderGateway:          *gateways.NewOrderGateway(orderDataSource),
		orderStatusGateway:    *gateways.NewOrderStatusGateway(orderStatusDataSource),
		messageBroker:         messageBroker,
	}
}

func (c *OrderController) Create(dto dtos.CreateOrderDTO) (dtos.OrderResponseDTO, error) {
	useCase := use_cases.NewCreateOrderUseCase(c.orderGateway, c.orderStatusGateway, c.messageBroker)
	order, err := useCase.Execute(dto.CustomerID, dto.Items)
	if err != nil {
		return dtos.OrderResponseDTO{}, err
	}
	return presenters.ToOrderResponse(order), nil
}

func (c *OrderController) FindAll(filter dtos.OrderFilterDTO) ([]dtos.OrderResponseDTO, error) {
	useCase := use_cases.NewFindAllOrdersUseCase(c.orderGateway)
	orders, err := useCase.Execute(filter)
	if err != nil {
		return nil, err
	}
	return presenters.ToOrderResponseList(orders), nil
}

func (c *OrderController) FindByID(id string) (dtos.OrderResponseDTO, error) {
	useCase := use_cases.NewFindOrderByIDUseCase(c.orderGateway)
	order, err := useCase.Execute(id)
	if err != nil {
		return dtos.OrderResponseDTO{}, err
	}
	return presenters.ToOrderResponse(order), nil
}

func (c *OrderController) Update(dto dtos.UpdateOrderDTO) (dtos.OrderResponseDTO, error) {
	useCase := use_cases.NewUpdateOrderUseCase(c.orderGateway, c.orderStatusGateway)
	order, err := useCase.Execute(dto)
	if err != nil {
		return dtos.OrderResponseDTO{}, err
	}
	return presenters.ToOrderResponse(order), nil
}

func (c *OrderController) UpdateStatus(dto dtos.UpdateOrderStatusDTO) (dtos.OrderResponseDTO, error) {
	useCase := use_cases.NewUpdateOrderStatusUseCase(&c.orderGateway, &c.orderStatusGateway)
	result, err := useCase.Execute(use_cases.UpdateOrderStatusDTO{
		OrderID: dto.OrderID,
		Status:  dto.Status,
	})
	if err != nil {
		return dtos.OrderResponseDTO{}, err
	}
	return presenters.ToOrderResponse(result.Order), nil
}

func (c *OrderController) Delete(id string) error {
	useCase := use_cases.NewDeleteOrderUseCase(c.orderGateway)
	return useCase.Execute(id)
}

func (c *OrderController) FindAllStatus() ([]dtos.OrderStatusResponseDTO, error) {
	useCase := use_cases.NewFindAllOrderStatusUseCase(c.orderStatusGateway)
	statuses, err := useCase.Execute()
	if err != nil {
		return nil, err
	}
	return presenters.ToOrderStatusResponseList(statuses), nil
}
