package service

import (
	"time"

	"github.com/fajarardiyanto/go-media-server/config"
	"github.com/fajarardiyanto/go-media-server/internal/model"
	"github.com/fajarardiyanto/go-media-server/internal/model/dto/request"
	"github.com/fajarardiyanto/go-media-server/internal/repository"
	"github.com/google/uuid"
)

type messageService struct{}

func NewMessageService() repository.MessageRepository {
	return &messageService{}
}

func (s *messageService) SendMessage(message request.RequestMessageModel) (*model.MessageModel, error) {
	data := model.MessageModel{
		Uuid:        uuid.NewString(),
		FromUser:    message.FromUser,
		ToUser:      message.ToUser,
		Content:     message.Content,
		MessageType: 0,
		CreatedAt:   time.Now(),
	}

	if err := config.GetDBConn().Orm().Debug().Model(&model.MessageModel{}).Create(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}
