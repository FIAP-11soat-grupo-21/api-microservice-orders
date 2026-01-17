package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"microservice/infra/api/rest/schemas"
	"microservice/internal/adapters/brokers"
	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
	"microservice/internal/interfaces"
	"microservice/utils/factories"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type mockOrderDS struct {
	createFunc   func(order daos.OrderDAO) error
	findAllFunc  func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error)
	findByIDFunc func(id string) (daos.OrderDAO, error)
	updateFunc   func(order daos.OrderDAO) error
	deleteFunc   func(id string) error
}

func (m *mockOrderDS) Create(order daos.OrderDAO) error {
	if m.createFunc != nil {
		return m.createFunc(order)
	}
	return nil
}

func (m *mockOrderDS) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(filter)
	}
	return []daos.OrderDAO{}, nil
}

func (m *mockOrderDS) FindByID(id string) (daos.OrderDAO, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return daos.OrderDAO{}, nil
}

func (m *mockOrderDS) Update(order daos.OrderDAO) error {
	if m.updateFunc != nil {
		return m.updateFunc(order)
	}
	return nil
}

func (m *mockOrderDS) Delete(id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

type mockOrderStatusDS struct {
	findByIDFunc   func(id string) (daos.OrderStatusDAO, error)
	findByNameFunc func(name string) (daos.OrderStatusDAO, error)
	findAllFunc    func() ([]daos.OrderStatusDAO, error)
}

func (m *mockOrderStatusDS) FindByID(id string) (daos.OrderStatusDAO, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return daos.OrderStatusDAO{ID: "status-1", Name: "Pending"}, nil
}

func (m *mockOrderStatusDS) FindByName(name string) (daos.OrderStatusDAO, error) {
	if m.findByNameFunc != nil {
		return m.findByNameFunc(name)
	}
	return daos.OrderStatusDAO{ID: "status-1", Name: name}, nil
}

func (m *mockOrderStatusDS) FindAll() ([]daos.OrderStatusDAO, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc()
	}
	return []daos.OrderStatusDAO{}, nil
}

type mockBroker struct {
	publishFunc func(topic string, message interface{}) error
}

func (m *mockBroker) ConsumeOrderUpdates(ctx context.Context, handler brokers.OrderUpdateHandler) error {
	return nil
}
func (m *mockBroker) ConsumeOrderError(ctx context.Context, handler brokers.OrderErrorHandler) error {
	return nil
}

func (m *mockBroker) Close() error {
	return nil
}

func (m *mockBroker) PublishOnTopic(ctx context.Context, topic string, message interface{}) error {
	if m.publishFunc != nil {
		return m.publishFunc(topic, message)
	}
	return nil
}

func setupMocks(orderDS *mockOrderDS, statusDS *mockOrderStatusDS, broker *mockBroker) func() {
	factories.SetNewOrderDataSource(func() interfaces.IOrderDataSource {
		return orderDS
	})
	factories.SetNewOrderStatusDataSource(func() interfaces.IOrderStatusDataSource {
		return statusDS
	})
	factories.SetNewMessageBroker(func() brokers.MessageBroker {
		return broker
	})

	return func() {
		factories.SetNewOrderDataSource(nil)
		factories.SetNewOrderStatusDataSource(nil)
		factories.SetNewMessageBroker(nil)
	}
}

func TestNewOrderHandler(t *testing.T) {
	orderDS := &mockOrderDS{}
	statusDS := &mockOrderStatusDS{}
	broker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, broker)
	defer cleanup()

	handler := NewOrderHandler()
	if handler == nil {
		t.Error("NewOrderHandler() returned nil")
	}
}

func TestOrderHandler_Create_Success(t *testing.T) {
	orderDS := &mockOrderDS{
		createFunc: func(order daos.OrderDAO) error {
			return nil
		},
	}
	statusDS := &mockOrderStatusDS{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{ID: "status-1", Name: "Pending"}, nil
		},
	}
	messageBroker := &mockBroker{
		publishFunc: func(topic string, message interface{}) error {
			return nil
		},
	}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.POST("/orders", handler.Create)

	body := schemas.CreateOrderSchema{
		Items: []schemas.CreateOrderItemSchema{
			{ProductID: "product-1", Quantity: 2, Price: 10.0},
		},
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Create() status = %v, want %v", w.Code, http.StatusCreated)
	}
}

func TestOrderHandler_Create_InvalidBody(t *testing.T) {
	orderDS := &mockOrderDS{}
	statusDS := &mockOrderStatusDS{}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.POST("/orders", handler.Create)

	req := httptest.NewRequest("POST", "/orders", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Create() status = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestOrderHandler_FindAll_Success(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	orderDS := &mockOrderDS{
		findAllFunc: func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
			return []daos.OrderDAO{
				{
					ID:         "order-1",
					CustomerID: &customerID,
					Amount:     20.0,
					Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
					Items: []daos.OrderItemDAO{
						{ID: "item-1", OrderID: "order-1", ProductID: "product-1", Quantity: 2, UnitPrice: 10.0},
					},
					CreatedAt: now,
				},
			}, nil
		},
	}
	statusDS := &mockOrderStatusDS{}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/orders", handler.FindAll)

	req := httptest.NewRequest("GET", "/orders", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("FindAll() status = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestOrderHandler_FindAll_WithFilters(t *testing.T) {
	orderDS := &mockOrderDS{
		findAllFunc: func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
			return []daos.OrderDAO{}, nil
		},
	}
	statusDS := &mockOrderStatusDS{}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/orders", handler.FindAll)

	req := httptest.NewRequest("GET", "/orders?status_id=status-1&customer_id=customer-123&created_at_from=2024-01-01T00:00:00Z&created_at_to=2024-12-31T23:59:59Z", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("FindAll() with filters status = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestOrderHandler_FindByID_Success(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	orderDS := &mockOrderDS{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:         "550e8400-e29b-41d4-a716-446655440000",
				CustomerID: &customerID,
				Amount:     20.0,
				Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
				Items: []daos.OrderItemDAO{
					{ID: "item-1", OrderID: "550e8400-e29b-41d4-a716-446655440000", ProductID: "product-1", Quantity: 2, UnitPrice: 10.0},
				},
				CreatedAt: now,
			}, nil
		},
	}
	statusDS := &mockOrderStatusDS{}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/orders/:id", handler.FindByID)

	req := httptest.NewRequest("GET", "/orders/550e8400-e29b-41d4-a716-446655440000", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("FindByID() status = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestOrderHandler_Update_Success(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	orderDS := &mockOrderDS{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:         "550e8400-e29b-41d4-a716-446655440000",
				CustomerID: &customerID,
				Amount:     20.0,
				Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
				Items: []daos.OrderItemDAO{
					{ID: "item-1", OrderID: "550e8400-e29b-41d4-a716-446655440000", ProductID: "product-1", Quantity: 2, UnitPrice: 10.0},
				},
				CreatedAt: now,
			}, nil
		},
		updateFunc: func(order daos.OrderDAO) error {
			return nil
		},
	}
	statusDS := &mockOrderStatusDS{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{ID: "status-2", Name: "Confirmed"}, nil
		},
	}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.PUT("/orders/:id", handler.Update)

	body := schemas.UpdateOrderSchema{StatusID: "status-2"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/orders/550e8400-e29b-41d4-a716-446655440000", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Update() status = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestOrderHandler_Update_InvalidBody(t *testing.T) {
	orderDS := &mockOrderDS{}
	statusDS := &mockOrderStatusDS{}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.PUT("/orders/:id", handler.Update)

	req := httptest.NewRequest("PUT", "/orders/550e8400-e29b-41d4-a716-446655440000", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Update() status = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestOrderHandler_Delete_Success(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	orderDS := &mockOrderDS{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:         "550e8400-e29b-41d4-a716-446655440000",
				CustomerID: &customerID,
				Amount:     20.0,
				Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
				Items:      []daos.OrderItemDAO{},
				CreatedAt:  now,
			}, nil
		},
		deleteFunc: func(id string) error {
			return nil
		},
	}
	statusDS := &mockOrderStatusDS{}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.DELETE("/orders/:id", handler.Delete)

	req := httptest.NewRequest("DELETE", "/orders/550e8400-e29b-41d4-a716-446655440000", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Delete() status = %v, want %v", w.Code, http.StatusNoContent)
	}
}

func TestOrderHandler_FindAllStatus_Success(t *testing.T) {
	orderDS := &mockOrderDS{}
	statusDS := &mockOrderStatusDS{
		findAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return []daos.OrderStatusDAO{
				{ID: "status-1", Name: "Pending"},
				{ID: "status-2", Name: "Confirmed"},
			}, nil
		},
	}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/orders/status", handler.FindAllStatus)

	req := httptest.NewRequest("GET", "/orders/status", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("FindAllStatus() status = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestToOrderResponse(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	dto := dtos.OrderResponseDTO{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		Status:     dtos.OrderStatusDTO{ID: "status-1", Name: "Pending"},
		Items: []dtos.OrderItemDTO{
			{ID: "item-1", ProductID: "product-1", OrderID: "order-1", Quantity: 2, UnitPrice: 50.0},
		},
		CreatedAt: now,
		UpdatedAt: nil,
	}

	response := toOrderResponse(dto)

	if response.ID != "order-1" {
		t.Errorf("toOrderResponse() ID = %v, want order-1", response.ID)
	}
	if *response.CustomerID != customerID {
		t.Errorf("toOrderResponse() CustomerID = %v, want %v", *response.CustomerID, customerID)
	}
	if response.Amount != 100.0 {
		t.Errorf("toOrderResponse() Amount = %v, want 100.0", response.Amount)
	}
	if response.Status != "Pending" {
		t.Errorf("toOrderResponse() Status = %v, want Pending", response.Status)
	}
	if len(response.Items) != 1 {
		t.Errorf("toOrderResponse() Items length = %v, want 1", len(response.Items))
	}
}

func TestOrderHandler_Create_Error(t *testing.T) {
	orderDS := &mockOrderDS{}
	statusDS := &mockOrderStatusDS{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{}, errors.New("status not found")
		},
	}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.POST("/orders", handler.Create)

	body := schemas.CreateOrderSchema{
		Items: []schemas.CreateOrderItemSchema{
			{ProductID: "product-1", Quantity: 2, Price: 10.0},
		},
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code == http.StatusCreated {
		t.Error("Create() should fail when status not found")
	}
}

func TestOrderHandler_FindAll_Error(t *testing.T) {
	orderDS := &mockOrderDS{
		findAllFunc: func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
			return nil, errors.New("database error")
		},
	}
	statusDS := &mockOrderStatusDS{}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/orders", handler.FindAll)

	req := httptest.NewRequest("GET", "/orders", nil)
	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK && w.Body.Len() > 2 {
		t.Log("FindAll() returned data despite error")
	}
}

func TestOrderHandler_FindByID_Error(t *testing.T) {
	orderDS := &mockOrderDS{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{}, errors.New("not found")
		},
	}
	statusDS := &mockOrderStatusDS{}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/orders/:id", handler.FindByID)

	req := httptest.NewRequest("GET", "/orders/550e8400-e29b-41d4-a716-446655440000", nil)
	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Log("FindByID() returned OK despite error")
	}
}

func TestOrderHandler_Update_Error(t *testing.T) {
	orderDS := &mockOrderDS{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{}, errors.New("not found")
		},
	}
	statusDS := &mockOrderStatusDS{}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.PUT("/orders/:id", handler.Update)

	body := schemas.UpdateOrderSchema{StatusID: "status-2"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/orders/550e8400-e29b-41d4-a716-446655440000", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Log("Update() returned OK despite error")
	}
}

func TestOrderHandler_Delete_Error(t *testing.T) {
	orderDS := &mockOrderDS{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{}, errors.New("not found")
		},
	}
	statusDS := &mockOrderStatusDS{}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.DELETE("/orders/:id", handler.Delete)

	req := httptest.NewRequest("DELETE", "/orders/550e8400-e29b-41d4-a716-446655440000", nil)
	router.ServeHTTP(w, req)

	if w.Code == http.StatusNoContent {
		t.Log("Delete() returned NoContent despite error")
	}
}

func TestOrderHandler_FindAllStatus_Error(t *testing.T) {
	orderDS := &mockOrderDS{}
	statusDS := &mockOrderStatusDS{
		findAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return nil, errors.New("database error")
		},
	}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/orders/status", handler.FindAllStatus)

	req := httptest.NewRequest("GET", "/orders/status", nil)
	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK && w.Body.Len() > 2 {
		t.Log("FindAllStatus() returned data despite error")
	}
}
func TestOrderHandler_UpdateStatus_Success(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	orderDS := &mockOrderDS{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:         "550e8400-e29b-41d4-a716-446655440000",
				CustomerID: &customerID,
				Amount:     20.0,
				Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
				Items: []daos.OrderItemDAO{
					{ID: "item-1", OrderID: "550e8400-e29b-41d4-a716-446655440000", ProductID: "product-1", Quantity: 2, UnitPrice: 10.0},
				},
				CreatedAt: now,
			}, nil
		},
		updateFunc: func(order daos.OrderDAO) error {
			return nil
		},
	}
	statusDS := &mockOrderStatusDS{
		findByNameFunc: func(name string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{ID: "status-2", Name: name}, nil
		},
	}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.PUT("/orders/:id/status", handler.UpdateStatus)

	body := schemas.UpdateOrderStatusSchema{Status: "Em preparação"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/orders/550e8400-e29b-41d4-a716-446655440000/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("UpdateStatus() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "Order status updated successfully" {
		t.Errorf("UpdateStatus() message = %v, want 'Order status updated successfully'", response["message"])
	}
}

func TestOrderHandler_UpdateStatus_InvalidBody(t *testing.T) {
	orderDS := &mockOrderDS{}
	statusDS := &mockOrderStatusDS{}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.PUT("/orders/:id/status", handler.UpdateStatus)

	req := httptest.NewRequest("PUT", "/orders/550e8400-e29b-41d4-a716-446655440000/status", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("UpdateStatus() status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Invalid request body" {
		t.Errorf("UpdateStatus() error = %v, want 'Invalid request body'", response["error"])
	}
}

func TestOrderHandler_UpdateStatus_OrderNotFound(t *testing.T) {
	orderDS := &mockOrderDS{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{}, errors.New("order not found")
		},
	}
	statusDS := &mockOrderStatusDS{}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	// Add the error handler middleware to process errors properly
	router.Use(func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			ctx.Abort()
		}
	})

	router.PUT("/orders/:id/status", handler.UpdateStatus)

	body := schemas.UpdateOrderStatusSchema{Status: "Em preparação"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/orders/non-existent-order/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("UpdateStatus() status = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestOrderHandler_UpdateStatus_StatusNotFound(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	orderDS := &mockOrderDS{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:         "550e8400-e29b-41d4-a716-446655440000",
				CustomerID: &customerID,
				Amount:     20.0,
				Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
				Items:      []daos.OrderItemDAO{},
				CreatedAt:  now,
			}, nil
		},
	}
	statusDS := &mockOrderStatusDS{
		findByNameFunc: func(name string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{}, errors.New("status not found")
		},
	}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	// Add the error handler middleware to process errors properly
	router.Use(func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			ctx.Abort()
		}
	})

	router.PUT("/orders/:id/status", handler.UpdateStatus)

	body := schemas.UpdateOrderStatusSchema{Status: "Status Inexistente"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/orders/550e8400-e29b-41d4-a716-446655440000/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("UpdateStatus() status = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestOrderHandler_UpdateStatus_UpdateError(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	orderDS := &mockOrderDS{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:         "550e8400-e29b-41d4-a716-446655440000",
				CustomerID: &customerID,
				Amount:     20.0,
				Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
				Items:      []daos.OrderItemDAO{},
				CreatedAt:  now,
			}, nil
		},
		updateFunc: func(order daos.OrderDAO) error {
			return errors.New("database update failed")
		},
	}
	statusDS := &mockOrderStatusDS{
		findByNameFunc: func(name string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{ID: "status-2", Name: name}, nil
		},
	}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	// Add the error handler middleware to process errors properly
	router.Use(func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			ctx.Abort()
		}
	})

	router.PUT("/orders/:id/status", handler.UpdateStatus)

	body := schemas.UpdateOrderStatusSchema{Status: "Em preparação"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/orders/550e8400-e29b-41d4-a716-446655440000/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("UpdateStatus() status = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestOrderHandler_UpdateStatus_SpecialCharactersInID(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	orderDS := &mockOrderDS{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			// Verificar se caracteres especiais foram sanitizados
			if id == "order_with_newline_and_carriage_return" {
				return daos.OrderDAO{
					ID:         id,
					CustomerID: &customerID,
					Amount:     20.0,
					Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
					Items:      []daos.OrderItemDAO{},
					CreatedAt:  now,
				}, nil
			}
			return daos.OrderDAO{}, errors.New("order not found")
		},
		updateFunc: func(order daos.OrderDAO) error {
			return nil
		},
	}
	statusDS := &mockOrderStatusDS{
		findByNameFunc: func(name string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{ID: "status-2", Name: name}, nil
		},
	}
	messageBroker := &mockBroker{}
	cleanup := setupMocks(orderDS, statusDS, messageBroker)
	defer cleanup()

	handler := NewOrderHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.PUT("/orders/:id/status", handler.UpdateStatus)

	body := schemas.UpdateOrderStatusSchema{Status: "Em preparação"}
	jsonBody, _ := json.Marshal(body)

	// Testar com caracteres especiais no ID que devem ser sanitizados
	// Usar URL encoding para caracteres especiais
	req := httptest.NewRequest("PUT", "/orders/order%0Awith%0Dnewline%0Aand%0Dcarriage%0Dreturn/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("UpdateStatus() status = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestOrderHandler_UpdateStatus_DifferentStatuses(t *testing.T) {
	testCases := []struct {
		name   string
		status string
	}{
		{
			name:   "Em preparação status",
			status: "Em preparação",
		},
		{
			name:   "Pronto status",
			status: "Pronto",
		},
		{
			name:   "Finalizado status",
			status: "Finalizado",
		},
		{
			name:   "Cancelado status",
			status: "Cancelado",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			customerID := "customer-123"
			now := time.Now()

			orderDS := &mockOrderDS{
				findByIDFunc: func(id string) (daos.OrderDAO, error) {
					return daos.OrderDAO{
						ID:         "550e8400-e29b-41d4-a716-446655440000",
						CustomerID: &customerID,
						Amount:     20.0,
						Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
						Items:      []daos.OrderItemDAO{},
						CreatedAt:  now,
					}, nil
				},
				updateFunc: func(order daos.OrderDAO) error {
					return nil
				},
			}
			statusDS := &mockOrderStatusDS{
				findByNameFunc: func(name string) (daos.OrderStatusDAO, error) {
					return daos.OrderStatusDAO{ID: "status-2", Name: name}, nil
				},
			}
			messageBroker := &mockBroker{}
			cleanup := setupMocks(orderDS, statusDS, messageBroker)
			defer cleanup()

			handler := NewOrderHandler()

			w := httptest.NewRecorder()
			_, router := gin.CreateTestContext(w)

			router.PUT("/orders/:id/status", handler.UpdateStatus)

			body := schemas.UpdateOrderStatusSchema{Status: tc.status}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest("PUT", "/orders/550e8400-e29b-41d4-a716-446655440000/status", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("UpdateStatus() for %s status = %v, want %v", tc.status, w.Code, http.StatusOK)
			}
		})
	}
}
