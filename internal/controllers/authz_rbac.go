package controllers

import (
	"fmt"

	"github.com/duvrdx/whoami/internal/schemas"
	"github.com/duvrdx/whoami/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthzRBACController struct {
	authzRBACService services.AuthzRBACService
}

func NewAuthzRBACController(authzRBACService services.AuthzRBACService) AuthzRBACController {
	return AuthzRBACController{authzRBACService: authzRBACService}
}

// RBACRole
func (controller AuthzRBACController) CreateRole(c echo.Context) error {
	var role schemas.RBACRoleCreate

	if err := c.Bind(&role); err != nil {
		return c.JSON(400, err)
	}

	createdRole, err := controller.authzRBACService.CreateRole(&role)
	if err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, createdRole)
}

func (controller AuthzRBACController) GetRole(c echo.Context) error {
	identifier := c.Param("identifier")

	role, err := controller.authzRBACService.GetRole(identifier)
	if err != nil {
		return c.JSON(404, err)
	}

	return c.JSON(200, role)
}

func (controller AuthzRBACController) GetRoles(c echo.Context) error {
	roles, err := controller.authzRBACService.GetRoles()
	if err != nil {
		return c.JSON(404, err)
	}

	return c.JSON(200, roles)
}

func (controller AuthzRBACController) UpdateRole(c echo.Context) error {
	identifier := c.Param("identifier")
	var role schemas.RBACRoleUpdate

	if err := c.Bind(&role); err != nil {
		return c.JSON(400, err)
	}

	updatedRole, err := controller.authzRBACService.UpdateRole(identifier, &role)
	if err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, updatedRole)
}

func (controller AuthzRBACController) DeleteRole(c echo.Context) error {
	identifier := c.Param("identifier")

	if err := controller.authzRBACService.DeleteRole(identifier); err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(204, "Role deleted successfully!")
}

// RBACPermission
func (controller AuthzRBACController) CreatePermission(c echo.Context) error {
	var permission schemas.RBACPermissionCreate

	if err := c.Bind(&permission); err != nil {
		return c.JSON(400, err)
	}

	createdPermission, err := controller.authzRBACService.CreatePermission(&permission)
	if err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, createdPermission)
}

func (controller AuthzRBACController) GetPermission(c echo.Context) error {
	identifier := c.Param("identifier")

	permission, err := controller.authzRBACService.GetPermission(identifier)
	if err != nil {
		return c.JSON(404, err)
	}

	return c.JSON(200, permission)
}

func (controller AuthzRBACController) GetPermissions(c echo.Context) error {

	permissions, err := controller.authzRBACService.GetPermissions()
	if err != nil {
		return c.JSON(404, err)
	}

	return c.JSON(200, permissions)
}

func (controller AuthzRBACController) UpdatePermission(c echo.Context) error {
	identifier := c.Param("identifier")
	var permission schemas.RBACPermissionUpdate

	if err := c.Bind(&permission); err != nil {
		return c.JSON(400, err)
	}

	updatedPermission, err := controller.authzRBACService.UpdatePermission(identifier, &permission)
	if err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, updatedPermission)
}

func (controller AuthzRBACController) DeletePermission(c echo.Context) error {
	identifier := c.Param("identifier")

	if err := controller.authzRBACService.DeletePermission(identifier); err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(204, "Permission deleted successfully!")
}

func (controller AuthzRBACController) GrantPermissionToRole(c echo.Context) error {
	roleIdentifier := c.Param("role")
	permissionIdentifier := c.Param("permission")

	if err := controller.authzRBACService.AddRoleToPermission(roleIdentifier, permissionIdentifier); err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, "Permission granted successfully!")
}

// RBACResourceType
func (controller AuthzRBACController) CreateResourceType(c echo.Context) error {
	var resourceType schemas.RBACResourceTypeCreate

	if err := c.Bind(&resourceType); err != nil {
		return c.JSON(400, err)
	}

	createdResourceType, err := controller.authzRBACService.CreateResourceType(&resourceType)
	if err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, createdResourceType)
}

func (controller AuthzRBACController) GetResourceType(c echo.Context) error {
	identifier := c.Param("identifier")

	resourceType, err := controller.authzRBACService.GetResourceType(identifier)
	if err != nil {
		return c.JSON(404, err)
	}

	return c.JSON(200, resourceType)
}

func (controller AuthzRBACController) GetResourceTypes(c echo.Context) error {
	resourceTypes, err := controller.authzRBACService.GetResourceTypes()
	if err != nil {
		return c.JSON(404, err)
	}

	return c.JSON(200, resourceTypes)
}

func (controller AuthzRBACController) UpdateResourceType(c echo.Context) error {
	identifier := c.Param("identifier")
	var resourceType schemas.RBACResourceTypeUpdate

	if err := c.Bind(&resourceType); err != nil {
		return c.JSON(400, err)
	}

	updatedResourceType, err := controller.authzRBACService.UpdateResourceType(identifier, &resourceType)
	if err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, updatedResourceType)
}

func (controller AuthzRBACController) DeleteResourceType(c echo.Context) error {
	identifier := c.Param("identifier")

	if err := controller.authzRBACService.DeleteResourceType(identifier); err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(204, "ResourceType deleted successfully!")
}

// Autorização
func (controller AuthzRBACController) AuthorizeByResourceType(c echo.Context) error {
	permissionIdentifier := c.QueryParam("permission")
	resourceTypeIdentifier := c.QueryParam("resource_type")

	userJWT := c.Get("user").(*schemas.UserResponse)

	authorized := controller.authzRBACService.AuthorizeUserByResourceType(userJWT.Identifier, permissionIdentifier, resourceTypeIdentifier)
	if !authorized {
		return c.JSON(403, "Forbidden")
	}

	return c.JSON(200, "Authorized")
}

func (controller AuthzRBACController) AuthorizeByResource(c echo.Context) error {
	permissionIdentifier := c.QueryParam("permission")
	resourceIdentifier := c.QueryParam("resource")
	userJWT := c.Get("user").(*schemas.UserResponse)

	fmt.Println("User JWT:", userJWT)

	authorized := controller.authzRBACService.AuthorizeUserByResource(userJWT.Identifier, permissionIdentifier, resourceIdentifier)
	if !authorized {
		return c.JSON(403, "Forbidden")
	}

	return c.JSON(200, "Authorized")
}

// Gerenciamento de Papéis e Permissões
func (controller AuthzRBACController) GrantRoleToUser(c echo.Context) error {
	roleIdentifier := c.QueryParam("role")
	userIdentifier := c.QueryParam("user")

	if err := controller.authzRBACService.GrantRoleToUser(roleIdentifier, userIdentifier); err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, "Role granted successfully!")
}

func (controller AuthzRBACController) RevokeRoleFromUser(c echo.Context) error {
	roleIdentifier := c.QueryParam("role")
	userIdentifier := c.QueryParam("user")

	if err := controller.authzRBACService.RevokeRoleFromUser(roleIdentifier, userIdentifier); err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, "Role revoked successfully!")
}

func (controller AuthzRBACController) ListGrantedRoles(c echo.Context) error {
	userIdentifier := c.QueryParam("user")

	roles, err := controller.authzRBACService.ListGrantedRoles(userIdentifier)
	if err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, roles)
}

func (controller AuthzRBACController) ListGrantedResources(c echo.Context) error {
	userIdentifier := c.QueryParam("user")
	permissionIdentifier := c.QueryParam("permission")

	resources, err := controller.authzRBACService.ListGrantedResources(userIdentifier, permissionIdentifier)
	if err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, resources)
}

func (controller AuthzRBACController) ListGrantedResourceTypes(c echo.Context) error {
	userIdentifier := c.QueryParam("user")
	permissionIdentifier := c.QueryParam("permission")

	resourceTypes, err := controller.authzRBACService.ListGrantedResourceTypes(userIdentifier, permissionIdentifier)
	if err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, resourceTypes)
}
