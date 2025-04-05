package routing

import (
	"github.com/duvrdx/whoami/internal/controllers"
	"github.com/duvrdx/whoami/internal/middlewares"
	"github.com/duvrdx/whoami/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Routing struct {
}

func (Routing Routing) GetRoutes() *echo.Echo {
	e := echo.New()

	e.Use(middlewares.LoggerMiddleware)
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middlewares.GetJWTMiddleware())

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

	// Authz RBAC routes
	authzRBACService := services.NewAuthzRBACService()
	authzRBACController := controllers.NewAuthzRBACController(authzRBACService)

	authz := e.Group("/authz")
	rbac := authz.Group("/rbac")
	rbac.POST("/role", authzRBACController.CreateRole)
	rbac.PUT("/role/:identifier", authzRBACController.UpdateRole)
	rbac.DELETE("/role/:identifier", authzRBACController.DeleteRole)
	rbac.GET("/role/:identifier", authzRBACController.GetRole)
	rbac.GET("/role", authzRBACController.GetRoles)
	rbac.POST("/role/grant", authzRBACController.GrantRoleToUser)
	rbac.POST("/role/revoke", authzRBACController.RevokeRoleFromUser)

	rbac.POST("/permission", authzRBACController.CreatePermission)
	rbac.PUT("/permission/:identifier", authzRBACController.UpdatePermission)
	rbac.DELETE("/permission/:identifier", authzRBACController.DeletePermission)
	rbac.GET("/permission/:identifier", authzRBACController.GetPermission)
	rbac.GET("/permission", authzRBACController.GetPermissions)
	rbac.POST("/permission/:permission/grant/:role", authzRBACController.GrantPermissionToRole)

	rbac.POST("/resourcetype", authzRBACController.CreateResourceType)
	rbac.PUT("/resourcetype/:identifier", authzRBACController.UpdateResourceType)
	rbac.DELETE("/resourcetype/:identifier", authzRBACController.DeleteResourceType)
	rbac.GET("/resourcetype/:identifier", authzRBACController.GetResourceType)
	rbac.GET("/resourcetype", authzRBACController.GetResourceTypes)

	rbac.POST("/authorize/resource", authzRBACController.AuthorizeByResource)
	rbac.POST("/authorize/resourcetype", authzRBACController.AuthorizeByResourceType)

	// rbac.GET("/role/granted", authzRBACController.ListGrantedRoles)
	// rbac.GET("/resource/granted", authzRBACController.ListGrantedResources)
	// rbac.GET("/resourcetype/granted", authzRBACController.ListGrantedResourceTypes)

	return e
}
