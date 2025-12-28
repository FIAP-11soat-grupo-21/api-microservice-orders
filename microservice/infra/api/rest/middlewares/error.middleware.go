package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"microservice/infra/api/rest/http_errors"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err

			errorHandled := http_errors.HandleDomainErrors(err, ctx)

			if !errorHandled {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}

			ctx.Abort()
		}
	}
}
