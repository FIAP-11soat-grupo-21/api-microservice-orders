package rest

import (
	"log"

	"github.com/gin-gonic/gin"

	"microservice/infra/api/rest/handlers"
	"microservice/infra/api/rest/middlewares"
	"microservice/infra/api/rest/routes"
	"microservice/infra/db/postgres"
	"microservice/infra/messaging"
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
		log.Printf("Warning: Failed to connect to RabbitMQ: %v", err)
		log.Println("The application will continue without message queue support")
	} else {
		paymentConsumer := messaging.NewPaymentEventConsumer()
		err = paymentConsumer.Start()
		if err != nil {
			log.Printf("Warning: Failed to start payment consumer: %v", err)
		}
	}

	ginRouter := NewRouter()
	if err := ginRouter.Run(":" + cfg.APIPort); err != nil {
		log.Fatalf("failed to start gin server: %v", err)
	}
}
