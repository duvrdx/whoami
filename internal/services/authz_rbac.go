package services

import (
	"github.com/duvrdx/whoami/internal/config"
	"github.com/duvrdx/whoami/internal/models"
	"github.com/duvrdx/whoami/internal/schemas"
	"github.com/duvrdx/whoami/internal/utils"
	"gorm.io/gorm"
)

type AuthzRBACService interface {
	CreateRole(role *schemas.RBACRoleCreate) (*schemas.RBACRoleResponse, error)
	GetRole(identifier string) (*schemas.RBACRoleResponse, error)
	GetRoles() ([]schemas.RBACRoleResponse, error)
	UpdateRole(identifier string, role *schemas.RBACRoleUpdate) (*schemas.RBACRoleResponse, error)
	DeleteRole(identifier string) error

	CreatePermission(permission *schemas.RBACPermissionCreate) (*schemas.RBACPermissionResponse, error)
	GetPermission(identifier string) (*schemas.RBACPermissionResponse, error)
	GetPermissions() ([]schemas.RBACPermissionResponse, error)
	UpdatePermission(identifier string, permission *schemas.RBACPermissionUpdate) (*schemas.RBACPermissionResponse, error)
	DeletePermission(identifier string) error
	AddRoleToPermission(roleIdentifier, permissionIdentifier string) error

	CreateResourceType(resourceType *schemas.RBACResourceTypeCreate) (*schemas.RBACResourceTypeResponse, error)
	GetResourceType(identifier string) (*schemas.RBACResourceTypeResponse, error)
	GetResourceTypes() ([]schemas.RBACResourceTypeResponse, error)
	UpdateResourceType(identifier string, resourceType *schemas.RBACResourceTypeUpdate) (*schemas.RBACResourceTypeResponse, error)
	DeleteResourceType(identifier string) error

	AuthorizeByResourceType(roleIdentifier, permissionIdentifier, resourceTypeIdentifier string) bool
	AuthorizeUserByResourceType(userIdentifier, permissionIdentifier, resourceTypeIdentifier string) bool
	AuthorizeByResource(roleIdentifier, permissionIdentifier, resourceIdentifier string) bool
	AuthorizeUserByResource(userIdentifier, permissionIdentifier, resourceIdentifier string) bool

	GrantRoleToUser(roleIdentifier, userIdentifier string) error
	RevokeRoleFromUser(roleIdentifier, userIdentifier string) error
	ListGrantedRoles(userIdentifier string) ([]string, error)
	ListGrantedResources(userIdentifier, permissionIdentifier string) ([]string, error)
	ListGrantedResourceTypes(userIdentifier, permissionIdentifier string) ([]string, error)
}

type authzRBACService struct {
	db *gorm.DB
}

// NewAuthzRBACService cria uma nova instância do serviço de RBAC
func NewAuthzRBACService() AuthzRBACService {
	return &authzRBACService{
		db: config.GetDB(),
	}
}

// RBACRole
func (s *authzRBACService) CreateRole(role *schemas.RBACRoleCreate) (*schemas.RBACRoleResponse, error) {
	roleModel := schemas.RBACRoleFromCreate(role)

	if err := s.db.Create(roleModel).Error; err != nil {
		return nil, err
	}

	returnRole := schemas.RBACRoleResponseFromModel(roleModel)

	return returnRole, nil
}

func (s *authzRBACService) GetRole(identifier string) (*schemas.RBACRoleResponse, error) {
	var role models.RBACRole

	if err := s.db.Where("identifier = ?", identifier).First(&role).Error; err != nil {
		return nil, err
	}

	returnRole := schemas.RBACRoleResponseFromModel(&role)

	return returnRole, nil
}

func (s *authzRBACService) GetRoles() ([]schemas.RBACRoleResponse, error) {
	var roles []models.RBACRole

	if err := s.db.Find(&roles).Error; err != nil {
		return nil, err
	}

	var returnRoles []schemas.RBACRoleResponse

	for _, role := range roles {
		returnRoles = append(returnRoles, *schemas.RBACRoleResponseFromModel(&role))
	}

	return returnRoles, nil
}

func (s *authzRBACService) UpdateRole(identifier string, role *schemas.RBACRoleUpdate) (*schemas.RBACRoleResponse, error) {
	var existing models.RBACRole

	if err := s.db.Where("identifier = ?", identifier).First(&existing).Error; err != nil {
		return nil, err
	}

	updateData := utils.MakeObjectWithoutNilFields(role)

	if len(updateData) == 0 {
		return schemas.RBACRoleResponseFromModel(&existing), nil
	}

	if err := s.db.Model(&existing).Updates(updateData).Error; err != nil {
		return nil, err
	}

	returnRole := schemas.RBACRoleResponseFromModel(&existing)
	return returnRole, nil
}

func (s *authzRBACService) DeleteRole(identifier string) error {
	var role models.RBACRole

	if err := s.db.Where("identifier = ?", identifier).First(&role).Error; err != nil {
		return err
	}

	if err := s.db.Delete(&role).Error; err != nil {
		return err
	}

	return nil
}

// RBACPermission

func (s *authzRBACService) GetPermission(identifier string) (*schemas.RBACPermissionResponse, error) {
	var permission models.RBACPermission

	if err := s.db.Where("identifier = ?", identifier).First(&permission).Error; err != nil {
		return nil, err
	}

	returnPermission := schemas.RBACPermissionResponseFromModel(&permission)

	return returnPermission, nil
}

func (s *authzRBACService) GetPermissions() ([]schemas.RBACPermissionResponse, error) {
	var permissions []models.RBACPermission

	if err := s.db.Find(&permissions).Error; err != nil {
		return nil, err
	}

	var returnPermissions []schemas.RBACPermissionResponse

	for _, permission := range permissions {
		returnPermissions = append(returnPermissions, *schemas.RBACPermissionResponseFromModel(&permission))
	}

	return returnPermissions, nil
}

func (s *authzRBACService) CreatePermission(permission *schemas.RBACPermissionCreate) (*schemas.RBACPermissionResponse, error) {
	permissionModel := schemas.RBACPermissionFromCreate(permission)

	if err := s.associateResources(permissionModel, permission.ResourceIdentifiers); err != nil {
		return nil, err
	}

	if err := s.db.Create(permissionModel).Error; err != nil {
		return nil, err
	}

	return schemas.RBACPermissionResponseFromModel(permissionModel), nil
}

func (s *authzRBACService) UpdatePermission(identifier string, permission *schemas.RBACPermissionUpdate) (*schemas.RBACPermissionResponse, error) {
	var existing models.RBACPermission
	if err := s.db.Where("identifier = ?", identifier).First(&existing).Error; err != nil {
		return nil, err
	}

	updateData := utils.MakeObjectWithoutNilFields(permission)

	if permission.ResourceIdentifiers != nil {
		if err := s.associateResources(&existing, permission.ResourceIdentifiers); err != nil {
			return nil, err
		}
	}

	if len(updateData) > 0 {
		if err := s.db.Model(&existing).Updates(updateData).Error; err != nil {
			return nil, err
		}
	}

	return schemas.RBACPermissionResponseFromModel(&existing), nil
}

func (s *authzRBACService) associateResources(permission *models.RBACPermission, resourceIdentifiers []string) error {
	var existingResources []models.RBACResourceIdentifier
	if err := s.db.Where("identifier IN ?", resourceIdentifiers).Find(&existingResources).Error; err != nil {
		return err
	}

	// Mapa dos existentes
	existingMap := make(map[string]bool)
	for _, res := range existingResources {
		existingMap[res.Identifier] = true
	}

	// Criar os que não existem
	var newResources []models.RBACResourceIdentifier
	for _, identifier := range resourceIdentifiers {
		if !existingMap[identifier] {
			newResources = append(newResources, models.RBACResourceIdentifier{
				Identifier: identifier,
			})
		}
	}

	if len(newResources) > 0 {
		if err := s.db.Create(&newResources).Error; err != nil {
			return err
		}
		existingResources = append(existingResources, newResources...)
	}

	// Associar todos ao permission
	return s.db.Model(permission).Association("ResourceIdentifiers").Replace(existingResources)
}

func (s *authzRBACService) DeletePermission(identifier string) error {
	var permission models.RBACPermission

	if err := s.db.Where("identifier = ?", identifier).First(&permission).Error; err != nil {
		return err
	}

	if err := s.db.Delete(&permission).Error; err != nil {
		return err
	}

	return nil
}

func (s *authzRBACService) AddRoleToPermission(roleIdentifier, permissionIdentifier string) error {
	var role models.RBACRole
	var permission models.RBACPermission

	// Verifica se o papel e a permissão existem
	if err := s.db.Where("identifier = ?", roleIdentifier).First(&role).Error; err != nil {
		return err
	}

	if err := s.db.Where("identifier = ?", permissionIdentifier).First(&permission).Error; err != nil {
		return err
	}

	// Associa o papel à permissão
	return s.db.Model(&permission).Association("AcceptedRoles").Append(&role)
}

// RBACResourceType
func (s *authzRBACService) CreateResourceType(resourceType *schemas.RBACResourceTypeCreate) (*schemas.RBACResourceTypeResponse, error) {
	resourceTypeModel := schemas.RBACResourceTypeFromCreate(resourceType)

	if err := s.db.Create(resourceTypeModel).Error; err != nil {
		return nil, err
	}

	returnResourceType := schemas.RBACResourceTypeResponseFromModel(resourceTypeModel)

	return returnResourceType, nil
}

func (s *authzRBACService) GetResourceType(identifier string) (*schemas.RBACResourceTypeResponse, error) {
	var resourceType models.RBACResourceType

	if err := s.db.Where("identifier = ?", identifier).First(&resourceType).Error; err != nil {
		return nil, err
	}

	returnResourceType := schemas.RBACResourceTypeResponseFromModel(&resourceType)

	return returnResourceType, nil
}

func (s *authzRBACService) GetResourceTypes() ([]schemas.RBACResourceTypeResponse, error) {
	var resourceTypes []models.RBACResourceType

	if err := s.db.Find(&resourceTypes).Error; err != nil {
		return nil, err
	}

	var returnResourceTypes []schemas.RBACResourceTypeResponse

	for _, resourceType := range resourceTypes {
		returnResourceTypes = append(returnResourceTypes, *schemas.RBACResourceTypeResponseFromModel(&resourceType))
	}

	return returnResourceTypes, nil
}

func (s *authzRBACService) UpdateResourceType(identifier string, resourceType *schemas.RBACResourceTypeUpdate) (*schemas.RBACResourceTypeResponse, error) {
	var existing models.RBACResourceType

	if err := s.db.Where("identifier = ?", identifier).First(&existing).Error; err != nil {
		return nil, err
	}

	updateData := utils.MakeObjectWithoutNilFields(resourceType)

	if len(updateData) == 0 {
		return schemas.RBACResourceTypeResponseFromModel(&existing), nil
	}

	if err := s.db.Model(&existing).Updates(updateData).Error; err != nil {
		return nil, err
	}

	returnResourceType := schemas.RBACResourceTypeResponseFromModel(&existing)
	return returnResourceType, nil
}

func (s *authzRBACService) DeleteResourceType(identifier string) error {
	var resourceType models.RBACResourceType

	if err := s.db.Where("identifier = ?", identifier).First(&resourceType).Error; err != nil {
		return err
	}

	if err := s.db.Delete(&resourceType).Error; err != nil {
		return err
	}

	return nil
}

// Autorização
func (s *authzRBACService) AuthorizeByResourceType(roleIdentifier, permissionIdentifier, resourceTypeIdentifier string) bool {
	var role models.RBACRole
	var permission models.RBACPermission
	var resourceType models.RBACResourceType

	// Verifica se o papel, a permissão e o tipo de recurso existem
	if err := s.db.Where("identifier = ?", roleIdentifier).First(&role).Error; err != nil {
		return false
	}

	if err := s.db.Where("identifier = ?", permissionIdentifier).First(&permission).Error; err != nil {
		return false
	}

	if err := s.db.Where("identifier = ?", resourceTypeIdentifier).First(&resourceType).Error; err != nil {
		return false
	}

	// Verifica se a permissão está associada ao tipo de recurso
	if permission.ResourceTypeID == nil || *permission.ResourceTypeID != resourceType.ID {
		return false
	}

	// Verifica se o papel tem a permissão
	exists := s.db.Model(&role).Where("id = ?", permission.ID).Association("Permissions").Count() > 0
	return exists

}

func (s *authzRBACService) AuthorizeUserByResourceType(userIdentifier, permissionIdentifier, resourceTypeIdentifier string) bool {
	user := models.User{}
	roles := []models.RBACRole{}

	if err := s.db.Where("identifier = ?", userIdentifier).First(&user).Error; err != nil {
		return false
	}

	if err := s.db.
		Joins("JOIN rbac_role_users ON rbac_role_users.rbac_role_id = rbac_roles.id").
		Where("rbac_role_users.user_id = ?", user.ID).
		Find(&roles).Error; err != nil {
		return false
	}

	for _, role := range roles {
		if s.AuthorizeUserByResourceType(role.Identifier, permissionIdentifier, resourceTypeIdentifier) {
			return true
		}
	}

	return false
}

func (s *authzRBACService) AuthorizeByResource(roleIdentifier, permissionIdentifier, resourceIdentifier string) bool {
	var role models.RBACRole
	var permission models.RBACPermission
	var resource models.RBACResourceIdentifier

	// Verifica se o papel, a permissão e o recurso existem
	if err := s.db.Where("identifier = ?", roleIdentifier).First(&role).Error; err != nil {
		return false
	}

	if err := s.db.Where("identifier = ?", permissionIdentifier).First(&permission).Error; err != nil {
		return false
	}

	if err := s.db.Where("identifier = ?", resourceIdentifier).First(&resource).Error; err != nil {
		return false
	}

	// Verifica se a permissão está associada ao recurso
	count := s.db.Model(&permission).Where("id = ?", resource.ID).Association("ResourceIdentifiers").Count()
	if count == 0 {
		return false
	}

	// Verifica se o papel tem a permissão
	count = s.db.Model(&role).Where("id = ?", permission.ID).Association("Permissions").Count()
	return count > 0
}

func (s *authzRBACService) AuthorizeUserByResource(userIdentifier, permissionIdentifier, resourceIdentifier string) bool {
	user := models.User{}
	roles := []models.RBACRole{}

	// Primeiro busca o usuário pelo identifier
	if err := s.db.Where("identifier = ?", userIdentifier).First(&user).Error; err != nil {
		return false
	}

	if err := s.db.
		Joins("JOIN rbac_role_users ON rbac_role_users.rbac_role_id = rbac_roles.id").
		Where("rbac_role_users.user_id = ?", user.ID).
		Find(&roles).Error; err != nil {
		return false
	}

	for _, role := range roles {
		if s.AuthorizeByResource(role.Identifier, permissionIdentifier, resourceIdentifier) {
			return true
		}
	}

	return false
}

// Gerenciamento de Papéis e Permissões
func (s *authzRBACService) GrantRoleToUser(roleIdentifier, userIdentifier string) error {
	var role models.RBACRole
	var user models.User

	// Verifica se o papel e o usuário existem
	if err := s.db.Where("identifier = ?", roleIdentifier).First(&role).Error; err != nil {
		return err
	}

	if err := s.db.Where("identifier = ?", userIdentifier).First(&user).Error; err != nil {
		return err
	}

	return s.db.Model(&role).Association("Users").Append(&user)
}

func (s *authzRBACService) RevokeRoleFromUser(roleIdentifier, userIdentifier string) error {
	var role models.RBACRole
	var user models.User

	// Verifica se o papel e o usuário existem
	if err := s.db.Where("identifier = ?", roleIdentifier).First(&role).Error; err != nil {
		return err
	}

	if err := s.db.Where("identifier = ?", userIdentifier).First(&user).Error; err != nil {
		return err
	}

	// Remove o papel do usuário
	return s.db.Model(&role).Association("Users").Delete(&user)
}

func (s *authzRBACService) ListGrantedRoles(userIdentifier string) ([]string, error) {
	var user models.User

	// Verifica se o usuário existe
	if err := s.db.Where("identifier = ?", userIdentifier).First(&user).Error; err != nil {
		return nil, err
	}

	// Lista os papéis associados ao usuário
	var roles []models.RBACRole
	if err := s.db.Model(&user).Association("Roles").Find(&roles); err != nil {
		return nil, err
	}

	roleIdentifiers := make([]string, len(roles))
	for i, role := range roles {
		roleIdentifiers[i] = role.Identifier
	}

	return roleIdentifiers, nil
}

func (s *authzRBACService) ListGrantedResources(userIdentifier, permissionIdentifier string) ([]string, error) {
	var user models.User
	var permission models.RBACPermission

	// Verifica se o usuário e a permissão existem
	if err := s.db.Where("identifier = ?", userIdentifier).First(&user).Error; err != nil {
		return nil, err
	}

	if err := s.db.Where("identifier = ?", permissionIdentifier).First(&permission).Error; err != nil {
		return nil, err
	}

	// Lista os recursos associados à permissão
	var resources []models.RBACResourceIdentifier
	if err := s.db.Model(&permission).Association("ResourceIdentifiers").Find(&resources); err != nil {
		return nil, err
	}

	resourceIdentifiers := make([]string, len(resources))
	for i, resource := range resources {
		resourceIdentifiers[i] = resource.Identifier
	}

	return resourceIdentifiers, nil
}

func (s *authzRBACService) ListGrantedResourceTypes(userIdentifier, permissionIdentifier string) ([]string, error) {
	var user models.User
	var permission models.RBACPermission

	// Verifica se o usuário e a permissão existem
	if err := s.db.Where("identifier = ?", userIdentifier).First(&user).Error; err != nil {
		return nil, err
	}

	if err := s.db.Where("identifier = ?", permissionIdentifier).First(&permission).Error; err != nil {
		return nil, err
	}

	// Lista os tipos de recurso associados à permissão
	var resourceTypes []models.RBACResourceType
	if err := s.db.Model(&permission).Association("ResourceType").Find(&resourceTypes); err != nil {
		return nil, err
	}

	resourceTypeIdentifiers := make([]string, len(resourceTypes))
	for i, resourceType := range resourceTypes {
		resourceTypeIdentifiers[i] = resourceType.Identifier
	}

	return resourceTypeIdentifiers, nil
}
