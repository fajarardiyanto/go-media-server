package model

import "time"

type UserTokenModel struct {
	ID       string
	Username string
	UserType int32
}

type UserModel struct {
	ID        string    `json:"id" gorm:"column:id"`
	Username  string    `json:"username" gorm:"column:username"`
	Password  string    `json:"-" gorm:"column:password"`
	UserType  int32     `json:"user_type" gorm:"column:user_type"`
	Status    bool      `json:"status" gorm:"column:status"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type UserReqModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserType int32  `json:"user_type"`
}

type UserResponseModel struct {
	User  UserModel `json:"user"`
	Token string    `json:"token"`
}

func (*UserModel) TableName() string {
	return "user"
}
