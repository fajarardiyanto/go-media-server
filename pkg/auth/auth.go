package auth

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/fajarardiyanto/go-media-server/config"
	"github.com/fajarardiyanto/go-media-server/internal/model"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func CreateToken(user model.UserTokenModel) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user"] = model.UserTokenModel{
		ID:       user.ID,
		Username: user.Username,
		UserType: user.UserType,
	}
	claims["exp"] = time.Now().Add(time.Hour * 20).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(model.GetConfig().ApiSecret))

}

func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(model.GetConfig().ApiSecret), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractToken(c *gin.Context) string {
	token, _ := c.GetQuery("token")
	if token != "" {
		return token
	}

	bearerToken := c.GetHeader("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenID(c *gin.Context) (*model.UserTokenModel, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(model.GetConfig().ApiSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		var user model.UserTokenModel

		userMarshal, err := json.Marshal(claims["user"].(map[string]interface{}))
		if err != nil {
			config.GetLogger().Error(err)
			return nil, err
		}

		if err = json.Unmarshal(userMarshal, &user); err != nil {
			config.GetLogger().Error(err)
			return nil, err
		}

		return &user, nil
	}
	return nil, nil
}
