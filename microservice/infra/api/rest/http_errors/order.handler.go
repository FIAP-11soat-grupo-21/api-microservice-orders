package http_errors

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"microservice/internal/domain/exceptions"
)

func HandleDomainErrors(err error, ctx *gin.Context) bool {
	switch e := err.(type) {
	case *exceptions.InvalidOrderDataException:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return true

	case *exceptions.OrderNotFoundException:
		ctx.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		return true

	case *exceptions.InvalidOrderItemData:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return true

	case *exceptions.AmountNotValidException:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return true

	case *exceptions.OrderStatusNotFoundException:
		ctx.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		return true
	}

	return false
}
