package auth

import (
	"context"
	"net/http"
	"time"

	"example.com/work-shop1/user"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Handler represents handler of login data
type Handler struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

var ctx = context.Background()

// Login User
func (h *Handler) Login(c echo.Context) error {
	password := c.FormValue("p")
	username := c.FormValue("u")

	user := user.User{}

	result := h.DB.Where("username=?", username).First(&user)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid UserName or Password",
		})
	}

	if user.Password != password {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid UserName or Password",
		})
	}

	token := uuid.New().String()

	// save  token to Radis
	err := h.RedisClient.Set(ctx, token, user.ID, 1*time.Hour).Err()
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})

}

//Validater To Validate user key after auth user
func (h *Handler) Validater(key string, c echo.Context) (bool, error) {

	uid, err := h.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return false, err
	}
	c.Set("uid", uid)
	return true, nil

}
