package models

import (
	"gorm.io/gorm"
)

type RBACRole struct {
	gorm.Model
	Identifier  string           `json:"identifier" gorm:"unique"`
	Name        string           `json:"name"`
	Description string           `json:"description" gorm:"default:''"`
	Permissions []RBACPermission `gorm:"many2many:rbac_role_permissions;"` // Relação many2many com RBACPermission
}

type RBACPermission struct {
	gorm.Model
	Identifier          string                   `json:"identifier" gorm:"unique"`
	Name                string                   `json:"name"`
	Description         string                   `json:"description" gorm:"default:''"`
	ResourceTypeID      *uint                    `json:"resourcetype_id"`                                 // Chave estrangeira para RBACResourceType
	ResourceIdentifiers []RBACResourceIdentifier `gorm:"many2many:rbac_permission_resource_identifiers;"` // Relação many2many com RBACResourceIdentifier

	ResourceType  *RBACResourceType `json:"resourcetype" gorm:"foreignKey:ResourceTypeID"` // Relação belongs to RBACResourceType
	AcceptedRoles []RBACRole        `gorm:"many2many:rbac_role_permissions;"`              // Relação many2many com RBACRole
}

type RBACResourceIdentifier struct {
	gorm.Model
	Identifier string `json:"identifier" gorm:"unique; uniqueIndex"`
}

type RBACResourceType struct {
	gorm.Model
	Identifier  string           `json:"identifier" gorm:"unique"`
	Name        string           `json:"name"`
	Description string           `json:"description" gorm:"default:''"`
	Permissions []RBACPermission `json:"permissions" gorm:"foreignKey:ResourceTypeID"` // Relação has many com RBACPermission
}
