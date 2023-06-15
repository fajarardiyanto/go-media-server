package handlers

import (
	"github.com/fajarardiyanto/go-media-server/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, model.Response{
		Error:   false,
		Message: "OK",
	})
}
