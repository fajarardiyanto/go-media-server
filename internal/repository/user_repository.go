package repository

import (
	"github.com/fajarardiyanto/go-media-server/internal/model"
)

type UserRepository interface {
	UserExist(username string) (*model.UserModel, error)
	Register(req model.UserModel) (*model.UserModel, error)
	GetUser() ([]model.UserModel, error)
	UpdateStatus(id string, status bool) error
	CheckUserLife(id string) bool
}
