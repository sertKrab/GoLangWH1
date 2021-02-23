package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	mylogger "example.com/work-shop1/logger"
	"github.com/go-redis/redis/v8"

	"example.com/work-shop1/auth"
	"example.com/work-shop1/post"
	"example.com/work-shop1/user"
)

var ctx = context.Background()

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:1433",
		Password: "GoLang", // no password set
		DB:       0,        // use default DB
	})

	// err := rdb.Set(ctx, "key", "value", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	// github.com/denisenkom/go-mssqldb
	// dsn := "sqlserver://sert11:1234567890@localhost:1433?database=social1"
	viper.SetDefault("dsn", "sqlserver://sert11:1234567890@localhost:1433?database=social1")
	viper.SetDefault("port", ":1323")
	viper.AutomaticEnv()
	db, err := gorm.Open(sqlserver.Open(viper.GetString("dsn")), &gorm.Config{})
	if err != nil {
		panic("error connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&post.Post{})
	db.AutoMigrate(&user.User{})
	// prepare handler
	h := &user.Handler{DB: db}
	p := &post.Handler{DB: db}
	authHandler := &auth.Handler{DB: db, RedisClient: rdb}

	// Echo instance
	e := echo.New()

	// Middleware
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(mylogger.Middleware(logger))

	// Routes
	//e.GET("/", hello)
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	e.GET("/hello", h.Hello)
	e.POST("/users", h.AddUser, middleware.KeyAuth(authHandler.Validater))
	e.GET("/users", h.ListUsers, middleware.KeyAuth(authHandler.Validater))
	e.GET("/users/:uid", h.GetUser, middleware.KeyAuth(authHandler.Validater))
	e.PUT("/users/:uid", h.UpdateUser, middleware.KeyAuth(authHandler.Validater))
	e.DELETE("/users/:uid", h.DeleteUser, middleware.KeyAuth(authHandler.Validater))

	e.POST("/users/:uid/posts", p.AddPost)
	e.GET("/users/:uid/posts", p.GetUserPosts)
	e.GET("/users/:uid/posts/:pid", p.GetUserPost)
	e.PUT("/users/:uid/posts/:pid", p.UpdateUserPost)
	e.DELETE("/users/:uid/posts/:pid", p.DeleteUserPost)

	e.POST("/login", authHandler.Login)
	// Start server
	//e.Logger.Fatal(e.Start(":1323"))
	// Start server
	go func() {
		if err := e.Start(viper.GetString("port")); err != nil {
			//e.Logger.Info("shutting down the server")
			logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
