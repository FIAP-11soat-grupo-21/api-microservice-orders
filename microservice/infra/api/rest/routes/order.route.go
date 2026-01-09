package routes

import (
	"github.com/gin-gonic/gin"

	"microservice/infra/api/rest/handlers"
)

func RegisterOrderRoutes(router *gin.RouterGroup) {
	handler := handlers.NewOrderHandler()

	router.POST("", handler.Create)
	router.GET("", handler.FindAll)
	router.GET("/:id", handler.FindByID)
	router.PUT("/:id", handler.Update)
	router.DELETE("/:id", handler.Delete)
}

func RegisterOrderStatusRoutes(router *gin.RouterGroup) {
	handler := handlers.NewOrderHandler()
	router.GET("/", handler.FindAllStatus)
}
