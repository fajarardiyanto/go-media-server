package middleware

import (
	"github.com/fajarardiyanto/go-media-server/internal/model"
	"github.com/fajarardiyanto/go-media-server/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := auth.TokenValid(c); err != nil {
			c.JSON(http.StatusUnauthorized, model.Response{
				Error:   true,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}
	}
}
