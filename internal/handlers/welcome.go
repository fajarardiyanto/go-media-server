package handlers

import (
	"net/http"

	"github.com/fajarardiyanto/go-media-server/internal/model/dto/response"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, response.Response{
		Error:   false,
		Message: "OK",
	})
}
