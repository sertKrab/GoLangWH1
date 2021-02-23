package logger

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

//Middleware return MiddlewareFunc
func Middleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			t := time.Now()
			id := c.Request().Header.Get("X-Request-ID")
			l := logger.With(zap.String("x-request-id", id))
			c.Set("logger", l)
			err := next(c)

			latency := time.Since(t)
			l.Info("request",
				zap.Int("status", c.Response().Status),
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.String("query", c.Request().URL.RawQuery),
				zap.String("ip", c.RealIP()),
				zap.String("user-agent", c.Request().UserAgent()),
				zap.Duration("latency", latency),
			)

			return err

		}
	}
}

//Extract zap logger  from context
func Extract(c echo.Context) *zap.Logger {
	l, ok := c.Get("logger").(*zap.Logger)
	if ok {
		return l
	}
	return nil
}
