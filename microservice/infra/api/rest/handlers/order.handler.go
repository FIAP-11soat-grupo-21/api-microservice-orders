package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"microservice/infra/api/rest/schemas"
	"microservice/infra/messaging"
	"microservice/internal/adapters/controllers"
	"microservice/internal/adapters/dtos"
	"microservice/utils/factories"
)

type OrderHandler struct {
	controller *controllers.OrderController
}

func NewOrderHandler() *OrderHandler {
	orderDataSource := factories.NewOrderDataSource()
	orderStatusDataSource := factories.NewOrderStatusDataSource()
	broker := messaging.GetBroker()

	if broker == nil {
		log.Println("Warning: Message broker not available, orders will be created without messaging")
	}

	controller := controllers.NewOrderController(orderDataSource, orderStatusDataSource, broker)

	return &OrderHandler{
		controller: controller,
	}
}

func (h *OrderHandler) Create(ctx *gin.Context) {
	var body schemas.CreateOrderSchema

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	items := make([]dtos.CreateOrderItemDTO, len(body.Items))
	for i, item := range body.Items {
		items[i] = dtos.CreateOrderItemDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	order, err := h.controller.Create(dtos.CreateOrderDTO{
		CustomerID: body.CustomerID,
		Items:      items,
	})

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, toOrderResponse(order))
}

func (h *OrderHandler) FindAll(ctx *gin.Context) {
	var filter dtos.OrderFilterDTO

	if createdAtFromStr := ctx.Query("created_at_from"); createdAtFromStr != "" {
		if t, err := time.Parse(time.RFC3339, createdAtFromStr); err == nil {
			filter.CreatedAtFrom = &t
		}
	}
	if createdAtToStr := ctx.Query("created_at_to"); createdAtToStr != "" {
		if t, err := time.Parse(time.RFC3339, createdAtToStr); err == nil {
			filter.CreatedAtTo = &t
		}
	}
	if statusID := ctx.Query("status_id"); statusID != "" {
		filter.StatusID = &statusID
	}
	if customerID := ctx.Query("customer_id"); customerID != "" {
		filter.CustomerID = &customerID
	}

	orders, err := h.controller.FindAll(filter)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	responses := make([]schemas.OrderResponseSchema, len(orders))
	for i, order := range orders {
		responses[i] = toOrderResponse(order)
	}

	ctx.JSON(http.StatusOK, responses)
}

func (h *OrderHandler) FindByID(ctx *gin.Context) {
	userInput := ctx.Param("id")
	orderID := strings.ReplaceAll(strings.ReplaceAll(userInput, "\n", "_"), "\r", "_")

	order, err := h.controller.FindByID(orderID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, toOrderResponse(order))
}

func (h *OrderHandler) Update(ctx *gin.Context) {
	userInput := ctx.Param("id")
	orderID := strings.ReplaceAll(strings.ReplaceAll(userInput, "\n", "_"), "\r", "_")

	var body schemas.UpdateOrderSchema
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	order, err := h.controller.Update(dtos.UpdateOrderDTO{
		ID:       orderID,
		StatusID: body.StatusID,
	})

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, toOrderResponse(order))
}

func (h *OrderHandler) Delete(ctx *gin.Context) {
	userInput := ctx.Param("id")
	orderID := strings.ReplaceAll(strings.ReplaceAll(userInput, "\n", "_"), "\r", "_")

	err := h.controller.Delete(orderID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *OrderHandler) FindAllStatus(ctx *gin.Context) {
	statuses, err := h.controller.FindAllStatus()
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	responses := make([]schemas.OrderStatusResponseSchema, len(statuses))
	for i, status := range statuses {
		responses[i] = schemas.OrderStatusResponseSchema{
			ID:   status.ID,
			Name: status.Name,
		}
	}

	ctx.JSON(http.StatusOK, responses)
}

func toOrderResponse(order dtos.OrderResponseDTO) schemas.OrderResponseSchema {
	items := make([]schemas.OrderItemResponseSchema, len(order.Items))
	for i, item := range order.Items {
		items[i] = schemas.OrderItemResponseSchema{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	return schemas.OrderResponseSchema{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Amount:     order.Amount,
		Status:     order.Status.Name,
		Items:      items,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}
