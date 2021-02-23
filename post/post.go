package post

import (
	"net/http"
	"strconv"

	"example.com/work-shop1/user"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Post structure
type Post struct {
	gorm.Model
	UserID  int
	User    user.User
	Content string
	Likes   int
}

// Handler represents handler of post data
type Handler struct {
	DB *gorm.DB
}

// AddPost api to Handler users Post Data
func (h Handler) AddPost(c echo.Context) error {

	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	post := Post{}
	err = c.Bind(&post)
	post.UserID = uid
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	result := h.DB.Create(&post)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, post)
}

//GetUserPosts to handler post of users
func (h Handler) GetUserPosts(c echo.Context) error {
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	posts := []Post{}
	//db.Where("name <> ?", "jinzhu").Find(&users)
	result := h.DB.Where("user_id = ?", uid).Find(&posts)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	} else if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "record not found",
		})
	}

	return c.JSON(http.StatusOK, posts)

}

// GetUserPost to handler getuserpost
func (h Handler) GetUserPost(c echo.Context) error {
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	post := Post{}

	result := h.DB.Where("user_id = ? AND id = ?", uid, pid).First(&post)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	} else if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "record not found",
		})
	}

	return c.JSON(http.StatusOK, post)
}

// UpdateUserPost to handler update post
func (h Handler) UpdateUserPost(c echo.Context) error {
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	post := Post{}

	result := h.DB.Where("user_id = ? AND id = ?", uid, pid).First(&post)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	} else if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "record not found",
		})
	}

	updatePost := Post{}

	err = c.Bind(&updatePost)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if updatePost.Content != "" {
		post.Content = updatePost.Content
	}

	if updatePost.Likes != 0 {
		post.Likes = updatePost.Likes
	}

	result = h.DB.Save(&post)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, post)
}

// DeleteUserPost to handler delete post
func (h Handler) DeleteUserPost(c echo.Context) error {
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	post := Post{}

	result := h.DB.Where("user_id = ? AND id = ?", uid, pid).First(&post)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	} else if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "record not found",
		})
	}

	result = h.DB.Delete(&post)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, post)

}
