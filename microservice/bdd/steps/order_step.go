package steps

import (
	"fmt"
	"microservice/infra/api/client"
	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
	"microservice/internal/adapters/gateways"
	"microservice/internal/domain/entities"
	"microservice/internal/use_cases"
	"microservice/mocks"

	"github.com/golang/mock/gomock"
)

type OrderHelper struct {
	Ctrl              *gomock.Controller
	MockOrderDS       *mocks.MockIOrderDataSource
	MockOrderStatusDS *mocks.MockIOrderStatusDataSource
	MockAC            *mocks.MockApiClient
	MockMB            *mocks.MockMessageBroker

	valid struct {
		CustomerID *string
		Items      []dtos.CreateOrderItemDTO
	}
	createdOrder entities.Order
	createError  error
}

func (oh *OrderHelper) TheOrderDataIsValid() error {
	customerID := "596d2388-d5ec-4b68-b6b8-ffa77273184a"
	oh.valid.CustomerID = &customerID
	oh.valid.Items = []dtos.CreateOrderItemDTO{
		{ProductID: "5653a03d-c6c8-481f-b723-f6c1c233d44c", Quantity: 2},
		{ProductID: "4c854e69-f0a0-4fe0-b350-f59701422d1f", Quantity: 1},
	}
	return nil
}

func (oh *OrderHelper) SendARequestToCreateANewOrder() error {
	// Setup mock expectations for FindByID (status lookup)
	statusDAO := daos.OrderStatusDAO{
		ID:   "56d3b3c3-1801-49cd-bae7-972c78082012",
		Name: "Recebido",
	}

	oh.MockOrderStatusDS.EXPECT().
		FindByID("56d3b3c3-1801-49cd-bae7-972c78082012").
		Return(statusDAO, nil).
		AnyTimes()

	// Setup mock expectations for API client (product lookup)
	oh.MockAC.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		DoAndReturn(func(path string, obj interface{}) error {
			if product, ok := obj.(*client.ProductResponseDTO); ok {
				product.ID = "product-1"
				product.Price = 10.0
				product.Active = true
			}
			return nil
		}).
		AnyTimes()

	// Setup mock expectation for Create
	oh.MockOrderDS.EXPECT().
		Create(gomock.Any()).
		Return(nil)

	// Setup mock expectation for PublishOnTopic
	oh.MockMB.EXPECT().
		PublishOnTopic(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	// Create gateways using the mocks
	orderGateway := gateways.NewOrderGateway(oh.MockOrderDS)
	statusGateway := gateways.NewOrderStatusGateway(oh.MockOrderStatusDS)

	// Create and execute the use case
	useCase := use_cases.NewCreateOrderUseCase(
		orderGateway,
		statusGateway,
		oh.MockMB,
		oh.MockAC,
	)

	oh.createdOrder, oh.createError = useCase.Execute(oh.valid.CustomerID, oh.valid.Items)

	return oh.createError
}

func (oh *OrderHelper) OrderShouldBeCreated() error {
	if oh.createError != nil {
		return fmt.Errorf("order creation failed: %w", oh.createError)
	}

	if oh.createdOrder.IsEmpty() {
		return fmt.Errorf("order was not created")
	}

	return nil
}
