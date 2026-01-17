package rest

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"microservice/infra/api/rest/handlers"
	"microservice/infra/api/rest/middlewares"
	"microservice/infra/api/rest/routes"
	"microservice/infra/db/postgres"
	"microservice/infra/db/postgres/data_source"
	"microservice/infra/messaging"
	"microservice/internal/adapters/consumers"
	"microservice/internal/adapters/gateways"
	"microservice/utils/config"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()

	ginRouter.Use(gin.Logger())
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(middlewares.ErrorHandlerMiddleware())

	healthHandler := handlers.NewHealthHandler()
	ginRouter.GET("/health", healthHandler.Health)

	v1Routes := ginRouter.Group("/v1")
	routes.RegisterOrderRoutes(v1Routes.Group("/orders"))
	routes.RegisterOrderStatusRoutes(v1Routes.Group("/orders/status"))

	return ginRouter
}

func Init() {
	cfg := config.LoadConfig()

	if cfg.IsProduction() {
		log.Printf("Running in production mode on [%s:%s]", cfg.APIHost, cfg.APIPort)
		gin.SetMode(gin.ReleaseMode)
	}

	postgres.Connect()

	if cfg.Database.RunMigrations {
		postgres.RunMigrations()
	}

	err := messaging.Connect()
	if err != nil {
		log.Printf("Warning: Failed to connect to message broker: %v", err)
		log.Println("The application will continue without message queue support")
	} else {
		log.Println("Message broker connected successfully")

		// Inicializar OrderUpdatesConsumer se broker estiver disponível
		broker := messaging.GetBroker()
		if broker != nil {
			ctx := context.Background()

			// Criar datasources e gateways necessários
			orderDataSource := data_source.NewGormOrderDataSource()
			orderStatusDataSource := data_source.NewGormOrderStatusDataSource()
			orderGateway := gateways.NewOrderGateway(orderDataSource)
			orderStatusGateway := gateways.NewOrderStatusGateway(orderStatusDataSource)

			// Criar consumer para atualizações de pedidos vindas do Kitchen Order
			orderUpdatesConsumer := consumers.NewOrderUpdatesConsumer(broker, orderGateway, orderStatusGateway)

			go func() {
				if err := orderUpdatesConsumer.Start(ctx); err != nil {
					log.Printf("Failed to start order updates consumer: %v", err)
				}
			}()

			orderErrorConsumer := consumers.NewOrderErrorConsumer(broker, orderGateway)

			go func() {
				if err := orderErrorConsumer.Start(ctx); err != nil {
					log.Printf("Failed to start order error consumer: %v", err)
				}
			}()
		}
	}

	ginRouter := NewRouter()
	if err := ginRouter.Run(":" + cfg.APIPort); err != nil {
		log.Fatalf("failed to start gin server: %v", err)
	}
}
