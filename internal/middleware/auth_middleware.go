package middleware

import (
	"net/http"

	"github.com/fajarardiyanto/go-media-server/internal/model/dto/response"
	"github.com/fajarardiyanto/go-media-server/pkg/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := auth.TokenValid(c); err != nil {
			c.JSON(http.StatusUnauthorized, response.Response{
				Error:   true,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}
	}
}
