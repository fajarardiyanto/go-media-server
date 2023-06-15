package repository

import (
	"github.com/fajarardiyanto/go-media-server/internal/model/dto/request"
)

type MessageRepository interface {
	SendMessage(message request.RequestMessageModel) error
}
