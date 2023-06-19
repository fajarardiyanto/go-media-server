package handlers

import (
	"net/http"
	"sync"

	"github.com/fajarardiyanto/go-media-server/config"
	"github.com/fajarardiyanto/go-media-server/internal/model"
	"github.com/fajarardiyanto/go-media-server/internal/model/dto/response"
	"github.com/fajarardiyanto/go-media-server/internal/repository"
	"github.com/fajarardiyanto/go-media-server/pkg/auth"
	"github.com/fajarardiyanto/go-media-server/util"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	sync.Mutex
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *userHandler {
	return &userHandler{repo: repo}
}

func (s *userHandler) RegisterHandler(c *gin.Context) {
	u := model.UserReqModel{}
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	if _, err := s.repo.UserExist(u.Username); err == nil {
		c.JSON(http.StatusInternalServerError, response.Response{
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
		c.JSON(http.StatusInternalServerError, response.Response{
			Error:   true,
			Message: "something went wrong while registering the user. please try again after sometime.",
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Error: false,
		Data:  res,
	})
}

func (s *userHandler) LoginHandler(c *gin.Context) {
	u := &model.UserReqModel{}
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	user, err := s.repo.UserExist(u.Username)
	if err != nil {
		config.GetLogger().Error(err)
		c.JSON(http.StatusInternalServerError, response.Response{
			Error:   true,
			Message: "Invalid username/password",
		})
		return
	}

	if err = util.VerifyPassword(user.Password, u.Password); err != nil {
		config.GetLogger().Error(err)
		c.JSON(http.StatusInternalServerError, response.Response{
			Error:   true,
			Message: "Invalid username/password",
		})
		return
	}

	userToken := model.UserTokenModel{
		ID:       user.ID,
		Username: user.Username,
		UserType: user.UserType,
	}

	token, err := auth.CreateToken(userToken)
	if err != nil {
		config.GetLogger().Error(err)
		c.JSON(http.StatusInternalServerError, response.Response{
			Error:   true,
			Message: "Something went wrong",
		})
		return
	}

	res := model.UserResponseModel{
		User:  *user,
		Token: token,
	}

	c.JSON(http.StatusOK, response.Response{
		Data: res,
	})
}
