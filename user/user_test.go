package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"example.com/work-shop1/post"
	"example.com/work-shop1/user"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var userHandler *user.Handler

func TestMain(m *testing.M) {
	dsn := "sqlserver://sert11:1234567890@localhost:1433?database=social1"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("error connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&post.Post{})
	db.AutoMigrate(&user.User{})
	// prepare handler
	userHandler = &user.Handler{
		DB: db,
	}

	os.Exit(m.Run())
}

func TestAddUser(t *testing.T) {

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Username": "test",
		"Password": "test_pass",
		"Name":     "test names",
		"Email":    "test@example",
	})

	given := string(givenBytes)
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(given))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := userHandler.AddUser(c); err != nil {
		t.Error(err)
		return
	}

	if rec.Code != http.StatusOK {
		t.Error("status code is not ok", rec.Code)
		return
	}

	returnUser := user.User{}

	err := json.Unmarshal(rec.Body.Bytes(), &returnUser)
	if err != nil {
		t.Error("can not Unmarshal respond", string(rec.Body.Bytes()))
		return
	}

	want := "test1"
	get := returnUser.Username

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

	want = "test names"
	get = returnUser.Name

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

	want = "test@example"
	get = returnUser.Email

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

}
