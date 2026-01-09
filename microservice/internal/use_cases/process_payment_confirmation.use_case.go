package use_cases

import (
	"fmt"
	"time"

	"microservice/internal/adapters/dtos"
	"microservice/internal/domain/entities"
	"microservice/internal/domain/exceptions"
	"microservice/internal/interfaces"
)

type ProcessPaymentConfirmationUseCase struct {
	orderGateway       interfaces.IOrderGateway
	orderStatusGateway interfaces.IOrderStatusGateway
}

func NewProcessPaymentConfirmationUseCase(
	orderGateway interfaces.IOrderGateway,
	orderStatusGateway interfaces.IOrderStatusGateway,
) *ProcessPaymentConfirmationUseCase {
	return &ProcessPaymentConfirmationUseCase{
		orderGateway:       orderGateway,
		orderStatusGateway: orderStatusGateway,
	}
}

type PaymentConfirmationDTO struct {
	OrderID       string
	PaymentID     string
	Status        string // "confirmed", "failed", "cancelled"
	Amount        float64
	PaymentMethod string
	ProcessedAt   time.Time
}

type PaymentConfirmationResult struct {
	Order               entities.Order
	StatusChanged       bool
	ShouldNotifyKitchen bool
	Message             string
}

func (uc *ProcessPaymentConfirmationUseCase) Execute(dto PaymentConfirmationDTO) (*PaymentConfirmationResult, error) {
	// 1. Validar dados de entrada
	if err := uc.validateInput(dto); err != nil {
		return nil, err
	}

	// 2. Buscar o pedido
	order, err := uc.orderGateway.FindByID(dto.OrderID)
	if err != nil {
		return nil, &exceptions.OrderNotFoundException{}
	}

	// 3. Verificar se o pedido pode ser atualizado
	if !uc.canUpdateOrderStatus(order, dto.Status) {
		return &PaymentConfirmationResult{
			Order:         *order,
			StatusChanged: false,
			Message:       fmt.Sprintf("Order %s cannot be updated from status %s", order.ID, order.Status.Name.Value()),
		}, nil
	}

	// 4. Processar baseado no status do pagamento
	switch dto.Status {
	case "confirmed":
		return uc.processConfirmedPayment(order, dto)
	case "failed", "cancelled":
		return uc.processFailedPayment(order, dto)
	default:
		return &PaymentConfirmationResult{
			Order:         *order,
			StatusChanged: false,
			Message:       fmt.Sprintf("Unknown payment status: %s", dto.Status),
		}, nil
	}
}

func (uc *ProcessPaymentConfirmationUseCase) processConfirmedPayment(order *entities.Order, dto PaymentConfirmationDTO) (*PaymentConfirmationResult, error) {
	// Buscar status "paid"
	paidStatus, err := uc.findOrCreateStatus("paid", "Paid")
	if err != nil {
		return nil, err
	}

	// Atualizar pedido
	updateDTO := dtos.UpdateOrderDTO{
		ID:       order.ID,
		StatusID: paidStatus.ID,
	}

	updatedOrder, err := uc.updateOrder(updateDTO)
	if err != nil {
		return nil, err
	}

	return &PaymentConfirmationResult{
		Order:               updatedOrder,
		StatusChanged:       true,
		ShouldNotifyKitchen: true,
		Message:             fmt.Sprintf("Order %s marked as paid", order.ID),
	}, nil
}

func (uc *ProcessPaymentConfirmationUseCase) processFailedPayment(order *entities.Order, dto PaymentConfirmationDTO) (*PaymentConfirmationResult, error) {
	// Buscar status "failed"
	failedStatus, err := uc.findOrCreateStatus("failed", "Failed")
	if err != nil {
		return nil, err
	}

	// Atualizar pedido
	updateDTO := dtos.UpdateOrderDTO{
		ID:       order.ID,
		StatusID: failedStatus.ID,
	}

	updatedOrder, err := uc.updateOrder(updateDTO)
	if err != nil {
		return nil, err
	}

	return &PaymentConfirmationResult{
		Order:               updatedOrder,
		StatusChanged:       true,
		ShouldNotifyKitchen: false,
		Message:             fmt.Sprintf("Order %s marked as failed", order.ID),
	}, nil
}

func (uc *ProcessPaymentConfirmationUseCase) validateInput(dto PaymentConfirmationDTO) error {
	if dto.OrderID == "" {
		return fmt.Errorf("order ID is required")
	}
	if dto.PaymentID == "" {
		return fmt.Errorf("payment ID is required")
	}
	if dto.Status == "" {
		return fmt.Errorf("payment status is required")
	}
	if dto.Amount <= 0 {
		return fmt.Errorf("payment amount must be positive")
	}
	return nil
}

func (uc *ProcessPaymentConfirmationUseCase) canUpdateOrderStatus(order *entities.Order, newStatus string) bool {
	currentStatus := order.Status.Name.Value()

	switch currentStatus {
	case "pending", "created":
		return newStatus == "confirmed" || newStatus == "failed" || newStatus == "cancelled"
	case "paid":
		return false // Pedido já pago não pode ser alterado
	case "failed", "cancelled":
		return false // Pedidos falhados/cancelados não podem ser alterados
	default:
		return false
	}
}

func (uc *ProcessPaymentConfirmationUseCase) findOrCreateStatus(statusID, statusName string) (*entities.OrderStatus, error) {
	status, err := uc.orderStatusGateway.FindByID(statusID)
	if err != nil {
		return entities.NewOrderStatus(statusID, statusName)
	}
	return status, nil
}

func (uc *ProcessPaymentConfirmationUseCase) updateOrder(dto dtos.UpdateOrderDTO) (entities.Order, error) {
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
