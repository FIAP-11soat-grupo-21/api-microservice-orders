package suites

import (
	"context"
	"microservice/bdd/steps"
	"microservice/mocks"
	"os"

	"github.com/cucumber/godog"
	"github.com/golang/mock/gomock"
)

type godogReporter struct{}

func (r *godogReporter) Errorf(_ string, _ ...interface{}) { /* forwarded to gomock, intentionally empty */
}
func (r *godogReporter) Fatalf(format string, _ ...interface{}) { panic("gomock fatal: " + format) }

// InitializeScenario registra os passos e controla o ciclo de vida dos mocks
func InitializeScenario(ctx *godog.ScenarioContext) {
	var helper *steps.OrderHelper

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		// Setup environment variables
		os.Setenv("API_GATEWAY_URL", "http://localhost:8081")
		os.Setenv("SNS_ORDER_CREATED_TOPIC_ARN", "arn:aws:sns:us-east-1:123456789012:order-created")

		ctrl := gomock.NewController(&godogReporter{})

		mockOrderDS := mocks.NewMockIOrderDataSource(ctrl)
		mockOrderStatusDS := mocks.NewMockIOrderStatusDataSource(ctrl)
		mockMB := mocks.NewMockMessageBroker(ctrl)
		mockAC := mocks.NewMockApiClient(ctrl)

		helper = &steps.OrderHelper{
			Ctrl:              ctrl,
			MockOrderDS:       mockOrderDS,
			MockOrderStatusDS: mockOrderStatusDS,
			MockMB:            mockMB,
			MockAC:            mockAC,
		}

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if helper != nil && helper.Ctrl != nil {
			helper.Ctrl.Finish()
		}
		// Cleanup environment variables
		os.Unsetenv("API_GATEWAY_URL")
		os.Unsetenv("SNS_ORDER_CREATED_TOPIC_ARN")
		return ctx, nil
	})

	ctx.Step(`^the order has valid data with customer ID and items$`, func() error { return helper.TheOrderDataIsValid() })
	ctx.Step(`^i send a request to create a new order$`, func() error { return helper.SendARequestToCreateANewOrder() })
	ctx.Step("^the order should be created successfully$", func() error { return helper.OrderShouldBeCreated() })
}
