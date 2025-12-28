package presenters

import (
	"microservice/internal/adapters/dtos"
	"microservice/internal/domain/entities"
)

func ToOrderResponse(order entities.Order) dtos.OrderResponseDTO {
	status := dtos.OrderStatusDTO{
		ID:   order.Status.ID,
		Name: order.Status.Name.Value(),
	}

	items := make([]dtos.OrderItemDTO, len(order.Items))
	for i, item := range order.Items {
		items[i] = dtos.OrderItemDTO{
			ID:        item.ID,
			ProductID: item.ProductID.Value(),
			OrderID:   item.OrderID,
			Quantity:  item.Quantity.Value(),
			UnitPrice: item.UnitPrice.Value(),
		}
	}

	return dtos.OrderResponseDTO{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Amount:     order.Amount.Value(),
		Status:     status,
		Items:      items,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func ToOrderResponseList(orders []entities.Order) []dtos.OrderResponseDTO {
	responses := make([]dtos.OrderResponseDTO, len(orders))
	for i, order := range orders {
		responses[i] = ToOrderResponse(order)
	}
	return responses
}

func ToOrderStatusResponse(status entities.OrderStatus) dtos.OrderStatusResponseDTO {
	return dtos.OrderStatusResponseDTO{
		ID:   status.ID,
		Name: status.Name.Value(),
	}
}

func ToOrderStatusResponseList(statuses []entities.OrderStatus) []dtos.OrderStatusResponseDTO {
	responses := make([]dtos.OrderStatusResponseDTO, len(statuses))
	for i, status := range statuses {
		responses[i] = ToOrderStatusResponse(status)
	}
	return responses
}
