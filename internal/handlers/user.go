package handlers

import (
	"github.com/fajarardiyanto/go-media-server/config"
	"github.com/fajarardiyanto/go-media-server/internal/model"
	"github.com/fajarardiyanto/go-media-server/internal/repository"
	"github.com/fajarardiyanto/go-media-server/pkg/auth"
	"github.com/fajarardiyanto/go-media-server/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

type UserHandler struct {
	sync.Mutex
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (s *UserHandler) RegisterHandler(c *gin.Context) {
	u := model.UserReqModel{}
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	if _, err := s.repo.UserExist(u.Username); err == nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Error:   true,
			Message: "username already exist!",
		})
		return
	}

	req := model.UserModel{
		Username: u.Username,
		Password: u.Password,
		UserType: u.UserType,
	}

	res, err := s.repo.Register(req)
	if err != nil {
		config.GetLogger().Error(err.Error())
		c.JSON(http.StatusInternalServerError, model.Response{
			Error:   true,
			Message: "something went wrong while registering the user. please try again after sometime.",
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Error: false,
		Data:  res,
	})
}

func (s *UserHandler) LoginHandler(c *gin.Context) {
	u := &model.UserReqModel{}
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	res, err := s.repo.UserExist(u.Username)
	if err != nil {
		config.GetLogger().Error(err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Error:   true,
			Message: "Invalid username/password",
		})
		return
	}

	if err = util.VerifyPassword(res.Password, u.Password); err != nil {
		config.GetLogger().Error(err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Error:   true,
			Message: "Invalid username/password",
		})
		return
	}

	userToken := model.UserTokenModel{
		ID:       res.ID,
		Username: res.Username,
		UserType: res.UserType,
	}

	token, err := auth.CreateToken(userToken)
	if err != nil {
		config.GetLogger().Error(err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Error:   true,
			Message: "Something went wrong",
		})
		return
	}

	response := model.UserResponseModel{
		User:  *res,
		Token: token,
	}

	c.JSON(http.StatusOK, model.Response{
		Data: response,
	})
}
