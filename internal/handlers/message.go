package handlers

import (
	"net/http"

	"github.com/fajarardiyanto/go-media-server/config"
	"github.com/fajarardiyanto/go-media-server/internal/model"
	"github.com/fajarardiyanto/go-media-server/internal/model/dto/request"
	"github.com/fajarardiyanto/go-media-server/internal/model/dto/response"
	"github.com/fajarardiyanto/go-media-server/internal/repository"
	"github.com/fajarardiyanto/go-media-server/pkg/broker"
	"github.com/gin-gonic/gin"
)

type messageHandler struct {
	messageRepository repository.MessageRepository
}

func NewMessageHandler(messageRepository repository.MessageRepository) *messageHandler {
	return &messageHandler{
		messageRepository: messageRepository,
	}
}

func (s *messageHandler) SendMessage(c *gin.Context) {
	message := request.RequestMessageModel{}
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	data, err := s.messageRepository.SendMessage(message)
	if err != nil {
		config.GetLogger().Error(err.Error())
		c.JSON(http.StatusInternalServerError, response.Response{
			Error:   true,
			Message: "something went wrong",
		})
		return
	}

	broker.OnMsg(model.GetConfig().Message, data)

	c.JSON(http.StatusOK, response.Response{
		Error: false,
		Data:  data,
	})
}
