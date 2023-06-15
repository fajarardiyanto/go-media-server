package service

import (
	"github.com/fajarardiyanto/go-media-server/config"
	"github.com/fajarardiyanto/go-media-server/internal/model"
	"github.com/fajarardiyanto/go-media-server/internal/repository"
	"github.com/fajarardiyanto/go-media-server/util"
	"github.com/google/uuid"
	"time"
)

type UserService struct{}

func NewUserService() repository.UserRepository {
	return &UserService{}
}

func (*UserService) UserExist(username string) (*model.UserModel, error) {
	var res model.UserModel
	if err := config.GetDBConn().Orm().Debug().Model(&model.UserModel{}).Where("username = ?", username).First(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}

func (*UserService) Register(req model.UserModel) (*model.UserModel, error) {
	pass, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.ID = uuid.NewString()
	req.Password = pass
	req.Status = false
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if err = config.GetDBConn().Orm().Debug().Model(&model.UserModel{}).Create(&req).Error; err != nil {
		return nil, err
	}

	return &req, nil
}

func (*UserService) GetUser() ([]model.UserModel, error) {
	var res []model.UserModel
	if err := config.GetDBConn().Orm().Debug().Model(&model.UserModel{}).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (*UserService) UpdateStatus(id string, status bool) error {
	if err := config.GetDBConn().Orm().Debug().Model(model.UserModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}

	return nil
}

func (*UserService) CheckUserLife(id string) bool {
	var res model.UserModel
	if err := config.GetDBConn().Orm().Debug().Model(&model.UserModel{}).Where("id = ?", id).First(&res).Error; err != nil {
		return false
	}

	return res.Status
}
