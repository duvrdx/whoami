package schemas

import (
	"github.com/duvrdx/whoami/internal/models"
)

// RBACRole schemas
type RBACRoleCreate struct {
	Identifier  string  `json:"identifier"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type RBACRoleUpdate struct {
	Identifier  *string `json:"identifier,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type RBACRoleResponse struct {
	ID          uint   `json:"id"`
	Identifier  string `json:"identifier"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func RBACRoleResponseFromModel(role *models.RBACRole) *RBACRoleResponse {
	return &RBACRoleResponse{
		ID:          role.ID,
		Identifier:  role.Identifier,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func RBACRoleFromCreate(role *RBACRoleCreate) *models.RBACRole {
	if role == nil {
		return nil
	}

	roleModel := &models.RBACRole{
		Identifier: role.Identifier,
		Name:       role.Name,
	}

	if role.Description != nil {
		roleModel.Description = *role.Description
	} else {
		roleModel.Description = ""
	}

	return roleModel
}

func RBACRoleFromUpdate(role *RBACRoleUpdate) *models.RBACRole {
	if role == nil {
		return nil
	}

	roleModel := &models.RBACRole{}

	if role.Identifier != nil {
		roleModel.Identifier = *role.Identifier
	}

	if role.Name != nil {
		roleModel.Name = *role.Name
	}

	if role.Description != nil {
		roleModel.Description = *role.Description
	}

	return roleModel
}

// RBACPermission schemas
type RBACPermissionCreate struct {
	Identifier          string   `json:"identifier"`
	Name                string   `json:"name"`
	Description         *string  `json:"description,omitempty"`
	ResourceTypeID      *uint    `json:"resourcetype_id,omitempty"`
	ResourceIdentifiers []string `json:"resource_identifiers,omitempty"`
}

type RBACPermissionUpdate struct {
	Identifier          *string  `json:"identifier,omitempty"`
	Name                *string  `json:"name,omitempty"`
	Description         *string  `json:"description,omitempty"`
	ResourceTypeID      *uint    `json:"resourcetype_id,omitempty"`
	ResourceIdentifiers []string `json:"resource_identifiers,omitempty"`
}

type RBACPermissionResponse struct {
	ID                  uint     `json:"id"`
	Identifier          string   `json:"identifier"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	ResourceTypeID      *uint    `json:"resourcetype_id,omitempty"`
	ResourceIdentifiers []string `json:"resource_identifiers,omitempty"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
}

func RBACPermissionResponseFromModel(permission *models.RBACPermission) *RBACPermissionResponse {
	resourceIdentifiers := make([]string, len(permission.ResourceIdentifiers))
	for i, identifier := range permission.ResourceIdentifiers {
		resourceIdentifiers[i] = identifier.Identifier
	}

	return &RBACPermissionResponse{
		ID:                  permission.ID,
		Identifier:          permission.Identifier,
		Name:                permission.Name,
		Description:         permission.Description,
		ResourceTypeID:      permission.ResourceTypeID,
		ResourceIdentifiers: resourceIdentifiers,
		CreatedAt:           permission.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:           permission.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func RBACPermissionFromCreate(permission *RBACPermissionCreate) *models.RBACPermission {
	if permission == nil {
		return nil
	}

	permissionModel := &models.RBACPermission{
		Identifier:     permission.Identifier,
		Name:           permission.Name,
		ResourceTypeID: permission.ResourceTypeID,
	}

	if permission.Description != nil {
		permissionModel.Description = *permission.Description
	} else {
		permissionModel.Description = ""
	}

	// Convert resource identifiers to models
	if len(permission.ResourceIdentifiers) > 0 {
		resourceIdentifiers := make([]models.RBACResourceIdentifier, len(permission.ResourceIdentifiers))
		for i, identifier := range permission.ResourceIdentifiers {
			resourceIdentifiers[i] = models.RBACResourceIdentifier{
				Identifier: identifier,
			}
		}
		permissionModel.ResourceIdentifiers = resourceIdentifiers
	}

	return permissionModel
}

func RBACPermissionFromUpdate(permission *RBACPermissionUpdate) *models.RBACPermission {
	if permission == nil {
		return nil
	}

	permissionModel := &models.RBACPermission{}

	if permission.Identifier != nil {
		permissionModel.Identifier = *permission.Identifier
	}

	if permission.Name != nil {
		permissionModel.Name = *permission.Name
	}

	if permission.Description != nil {
		permissionModel.Description = *permission.Description
	}

	if permission.ResourceTypeID != nil {
		permissionModel.ResourceTypeID = permission.ResourceTypeID
	}

	// Convert resource identifiers to models
	if len(permission.ResourceIdentifiers) > 0 {
		resourceIdentifiers := make([]models.RBACResourceIdentifier, len(permission.ResourceIdentifiers))
		for i, identifier := range permission.ResourceIdentifiers {
			resourceIdentifiers[i] = models.RBACResourceIdentifier{
				Identifier: identifier,
			}
		}
		permissionModel.ResourceIdentifiers = resourceIdentifiers
	}

	return permissionModel
}

// RBACResourceType schemas
type RBACResourceTypeCreate struct {
	Identifier  string  `json:"identifier"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type RBACResourceTypeUpdate struct {
	Identifier  *string `json:"identifier,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type RBACResourceTypeResponse struct {
	ID          uint   `json:"id"`
	Identifier  string `json:"identifier"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func RBACResourceTypeResponseFromModel(resourceType *models.RBACResourceType) *RBACResourceTypeResponse {
	return &RBACResourceTypeResponse{
		ID:          resourceType.ID,
		Identifier:  resourceType.Identifier,
		Name:        resourceType.Name,
		Description: resourceType.Description,
		CreatedAt:   resourceType.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   resourceType.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func RBACResourceTypeFromCreate(resourceType *RBACResourceTypeCreate) *models.RBACResourceType {
	if resourceType == nil {
		return nil
	}

	resourceTypeModel := &models.RBACResourceType{
		Identifier: resourceType.Identifier,
		Name:       resourceType.Name,
	}

	if resourceType.Description != nil {
		resourceTypeModel.Description = *resourceType.Description
	} else {
		resourceTypeModel.Description = ""
	}

	return resourceTypeModel
}

func RBACResourceTypeFromUpdate(resourceType *RBACResourceTypeUpdate) *models.RBACResourceType {
	if resourceType == nil {
		return nil
	}

	resourceTypeModel := &models.RBACResourceType{}

	if resourceType.Identifier != nil {
		resourceTypeModel.Identifier = *resourceType.Identifier
	}

	if resourceType.Name != nil {
		resourceTypeModel.Name = *resourceType.Name
	}

	if resourceType.Description != nil {
		resourceTypeModel.Description = *resourceType.Description
	}

	return resourceTypeModel
}
