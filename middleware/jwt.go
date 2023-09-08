package middleware

import (
	"backend_myblog/model"
	"backend_myblog/security"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func JWTMiddleware() echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:     &model.JwtCustomClaim{},
		SigningKey: []byte(security.Secret_key),
	}

	return middleware.JWTWithConfig(config)
}