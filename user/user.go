package user

import (
	"net/http"
	"strconv"

	"example.com/work-shop1/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// User entity for authentication
type User struct {
	gorm.Model
	Username string
	Password string
	Name     string
	Email    string
}

// Handler represents handler of user data
type Handler struct {
	DB *gorm.DB
}

// Hello for test api Respont
func (h *Handler) Hello(c echo.Context) error {
	return c.JSON(http.StatusOK, "User Hello ok")
}

// AddUser api handler add user Request
func (h *Handler) AddUser(c echo.Context) error {
	user := User{}
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	result := h.DB.Create(&user)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)

}

// ListUsers Handler to getlistUser Request
func (h *Handler) ListUsers(c echo.Context) error {
	users := []User{}

	l := logger.Extract(c)
	l.Info("listing users")
	result := h.DB.Find(&users)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	}

	l.Info("Returing users")

	return c.JSON(http.StatusOK, users)

}

// GetUser Handler to GetUsers Request
func (h *Handler) GetUser(c echo.Context) error {
	uid, err := strconv.Atoi(c.Param("uid"))
	l := logger.Extract(c)
	l.Info("getting user", zap.Int("uid", uid))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	user := User{}

	result := h.DB.Find(&user, uid)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	} else if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "record not found",
		})
	}

	return c.JSON(http.StatusOK, user)

}

// UpdateUser api handler Update user Request
func (h *Handler) UpdateUser(c echo.Context) error {

	loginUID := c.Get("uid").(string)
	if loginUID != c.Param("uid") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized user",
		})
	}

	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	user := User{}

	result := h.DB.Find(&user, uid)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	} else if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "record not found",
		})
	}

	updateUser := User{}

	err = c.Bind(&updateUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if updateUser.Name != "" {
		user.Name = updateUser.Name
	}

	if updateUser.Password != "" {
		user.Password = updateUser.Password
	}

	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}

	result = h.DB.Save(&user)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)

}

// DeleteUser api handler delete user Request
func (h *Handler) DeleteUser(c echo.Context) error {
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	user := User{}

	result := h.DB.Find(&user, uid)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	} else if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "record not found",
		})
	}

	result = h.DB.Delete(&user)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)

}
