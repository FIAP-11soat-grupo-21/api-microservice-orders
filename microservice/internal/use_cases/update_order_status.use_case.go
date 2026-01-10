package use_cases

import (
	"fmt"
	"log"

	"microservice/internal/domain/entities"
	"microservice/internal/interfaces"
)

type UpdateOrderStatusUseCase struct {
	orderGateway       interfaces.IOrderGateway
	orderStatusGateway interfaces.IOrderStatusGateway
}

type UpdateOrderStatusDTO struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type UpdateOrderStatusResult struct {
	Order   entities.Order `json:"order"`
	Message string         `json:"message"`
}

func NewUpdateOrderStatusUseCase(orderGateway interfaces.IOrderGateway, orderStatusGateway interfaces.IOrderStatusGateway) *UpdateOrderStatusUseCase {
	return &UpdateOrderStatusUseCase{
		orderGateway:       orderGateway,
		orderStatusGateway: orderStatusGateway,
	}
}

func (uc *UpdateOrderStatusUseCase) Execute(dto UpdateOrderStatusDTO) (*UpdateOrderStatusResult, error) {
	log.Printf("Updating order %s status to: %s", dto.OrderID, dto.Status)

	// Buscar o pedido
	order, err := uc.orderGateway.FindByID(dto.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order %s: %w", dto.OrderID, err)
	}

	// Mapear o status da cozinha para o status do pedido
	orderStatusName := uc.mapKitchenStatusToOrderStatus(dto.Status)

	// Buscar o status pelo nome
	newStatus, err := uc.orderStatusGateway.FindByName(orderStatusName)
	if err != nil {
		return nil, fmt.Errorf("failed to find order status '%s': %w", orderStatusName, err)
	}

	// Atualizar o status do pedido
	order.Status = *newStatus

	// Salvar as alterações
	err = uc.orderGateway.Update(*order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order %s: %w", dto.OrderID, err)
	}

	result := &UpdateOrderStatusResult{
		Order:   *order,
		Message: fmt.Sprintf("Order %s status updated to %s", dto.OrderID, orderStatusName),
	}

	log.Printf("Order %s status successfully updated to: %s", dto.OrderID, orderStatusName)
	return result, nil
}

func (uc *UpdateOrderStatusUseCase) mapKitchenStatusToOrderStatus(kitchenStatus string) string {
	statusMap := map[string]string{
		"Em preparação": "Em preparação",
		"Pronto":        "Pronto",
		"Finalizado":    "Finalizado",
	}

	if orderStatus, exists := statusMap[kitchenStatus]; exists {
		return orderStatus
	}

	return kitchenStatus
}
