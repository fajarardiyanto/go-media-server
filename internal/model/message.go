package model

import "time"

type MessageModel struct {
	ID          int32     `gorm:"column:id"`
	Uuid        string    `gorm:"column:uuid"`
	FromUser    string    `gorm:"column:from_user"`
	ToUser      string    `gorm:"column:to_user"`
	Content     string    `gorm:"column:content"`
	MessageType int32     `gorm:"column:message_type"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (*MessageModel) TableName() string {
	return "messages"
}
