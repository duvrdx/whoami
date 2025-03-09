package middlewares

import "github.com/labstack/echo/v4/middleware"

var LoggerMiddleware = middleware.LoggerWithConfig(middleware.LoggerConfig{
	Format: "method=${method}, uri=${uri}, status=${status}\n",
})
