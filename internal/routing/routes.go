package routing

import (
	"github.com/duvrdx/whoami/internal/controllers"
	"github.com/duvrdx/whoami/internal/middlewares"
	"github.com/duvrdx/whoami/internal/services"
	"github.com/labstack/echo/v4"
)

type Routing struct {
}

func (Routing Routing) GetRoutes() *echo.Echo {
	e := echo.New()

	e.Use(middlewares.JWTMiddleware)
	e.Use(middlewares.LoggerMiddleware)

	// Controllers and Services definitions
	authService := services.NewAuthService()
	authController := controllers.NewAuthController(authService)

	// OAuth2 routes
	oauth := e.Group("/o")
	oauth.POST("/token", authController.Token)
	oauth.DELETE("/token/:identifier", authController.RevokeToken)
	oauth.POST("/token/authorize", authController.Authorize)
	oauth.POST("/token/refresh", authController.RefreshToken)

	// Auth routes
	auth := e.Group("/auth")
	auth.POST("/user", authController.Register)
	auth.PUT("/user/:identifier", authController.UpdateUser)
	auth.DELETE("/user/:identifier", authController.DeleteUser)
	auth.GET("/user/:identifier", authController.GetUser)
	auth.GET("/user", authController.GetUsers)

	auth.POST("/client", authController.CreateClient)
	auth.PUT("/client/:identifier", authController.UpdateClient)
	auth.DELETE("/client/:identifier", authController.DeleteClient)
	auth.GET("/client/:identifier", authController.GetClient)
	auth.GET("/client", authController.GetClients)

	return e
}
