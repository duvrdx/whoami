package middlewares

import (
	"strings"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("your-secret-key")

var config = echojwt.Config{
	SigningKey: jwtSecret,
	Skipper: func(c echo.Context) bool {
		return strings.HasPrefix(c.Path(), "/oauth")
	},
}

var JWTMiddleware = echojwt.WithConfig(config)
