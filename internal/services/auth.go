package services

import (
	"time"

	"github.com/duvrdx/whoami/internal/config"
	"github.com/duvrdx/whoami/internal/models"
	"github.com/duvrdx/whoami/internal/schemas"
	"github.com/duvrdx/whoami/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AuthService interface
type AuthService interface {
	CreateUser(user *schemas.UserCreate) (*schemas.UserResponse, error)
	GetUser(identifier string) (*schemas.UserResponse, error)
	GetUsers() ([]schemas.UserResponse, error)
	UpdateUser(identifier string, user *schemas.UserUpdate) (*schemas.UserResponse, error)
	DeleteUser(identifier string) error
	CompareUserPassword(identifier, hashedPassword string) bool

	CreateClient(client *schemas.ClientCreate) (*schemas.ClientResponse, error)
	GetClient(identifier string) (*schemas.ClientResponse, error)
	GetClients() ([]schemas.ClientResponse, error)
	UpdateClient(identifier string, client *schemas.ClientUpdate) (*schemas.ClientResponse, error)
	DeleteClient(identifier string) error
	VerifyClient(identifier, secret string) bool

	CreateGroup(group *schemas.GroupCreate) (*schemas.GroupResponse, error)
	GetGroup(identifier string) (*schemas.GroupResponse, error)
	UpdateGroup(identifier string, group *schemas.GroupUpdate) (*schemas.GroupResponse, error)
	DeleteGroup(identifier string) error

	CreateToken(token *schemas.TokenCreate) (*schemas.TokenResponse, error)
	GetToken(identifier string) (*schemas.TokenResponse, error)
	GetTokenByRefreshToken(refreshToken string) (*models.Token, error)
	GetTokenByAccessToken(accessToken string) (*schemas.TokenResponse, error)

	RevokeToken(identifier string) error
	Authorize(accessToken string) bool
}

// AuthService implementation
type authService struct {
	db *gorm.DB
}

// NewAuthService creates a new auth service
func NewAuthService() AuthService {
	return &authService{
		db: config.GetDB(),
	}
}

func (s *authService) CreateUser(user *schemas.UserCreate) (*schemas.UserResponse, error) {

	userModel := schemas.UserFromCreate(user)

	if err := s.db.Create(userModel).Error; err != nil {
		return nil, err
	}

	returnUser := schemas.UserResponseFromModel(userModel)

	return returnUser, nil
}

func (s *authService) GetUser(identifier string) (*schemas.UserResponse, error) {
	var user models.User

	if err := s.db.Where("identifier = ?", identifier).First(&user).Error; err != nil {
		return nil, err
	}

	returnUser := schemas.UserResponseFromModel(&user)

	return returnUser, nil
}

func (s *authService) GetUsers() ([]schemas.UserResponse, error) {
	var users []models.User

	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	var returnUsers []schemas.UserResponse

	for _, user := range users {
		returnUsers = append(returnUsers, *schemas.UserResponseFromModel(&user))
	}

	return returnUsers, nil
}

func (s *authService) UpdateUser(identifier string, user *schemas.UserUpdate) (*schemas.UserResponse, error) {
	var existing models.User

	// Verifica se o usuário existe
	if err := s.db.Where("identifier = ?", identifier).First(&existing).Error; err != nil {
		return nil, err
	}

	updateData := utils.MakeObjectWithoutNilFields(user)

	if len(updateData) == 0 {
		return schemas.UserResponseFromModel(&existing), nil
	}

	if err := s.db.Model(&existing).Updates(updateData).Error; err != nil {
		return nil, err
	}

	returnUser := schemas.UserResponseFromModel(&existing)
	return returnUser, nil
}

func (s *authService) DeleteUser(identifier string) error {
	var user models.User

	if err := s.db.Where("identifier = ?", identifier).First(&user).Error; err != nil {
		return err
	}

	if err := s.db.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *authService) CompareUserPassword(identifier, password string) bool {
	var user models.User

	if err := s.db.Where("identifier = ?", identifier).First(&user).Error; err != nil {
		return false
	}

	return utils.CheckPassword(password, user.Password)
}

func (s *authService) CreateClient(client *schemas.ClientCreate) (*schemas.ClientResponse, error) {
	clientModel := schemas.ClientFromCreate(client)

	if err := s.db.Create(clientModel).Error; err != nil {
		return nil, err
	}

	returnClient := schemas.ClientResponseFromModel(clientModel)

	return returnClient, nil
}

func (s *authService) GetClient(identifier string) (*schemas.ClientResponse, error) {
	var client models.Client

	if err := s.db.Where("identifier = ?", identifier).First(&client).Error; err != nil {
		return nil, err
	}

	returnClient := schemas.ClientResponseFromModel(&client)

	return returnClient, nil
}

func (s *authService) GetClients() ([]schemas.ClientResponse, error) {
	var clients []models.Client

	if err := s.db.Find(&clients).Error; err != nil {
		return nil, err
	}

	var returnClients []schemas.ClientResponse

	for _, client := range clients {
		returnClients = append(returnClients, *schemas.ClientResponseFromModel(&client))
	}

	return returnClients, nil
}

func (s *authService) UpdateClient(identifier string, client *schemas.ClientUpdate) (*schemas.ClientResponse, error) {
	var existing models.Client

	if err := s.db.Where("identifier = ?", identifier).First(&existing).Error; err != nil {
		return nil, err
	}

	updateData := utils.MakeObjectWithoutNilFields(client)

	if len(updateData) == 0 {
		return schemas.ClientResponseFromModel(&existing), nil
	}

	if err := s.db.Model(&existing).Updates(updateData).Error; err != nil {
		return nil, err
	}

	returnClient := schemas.ClientResponseFromModel(&existing)
	return returnClient, nil
}

func (s *authService) DeleteClient(identifier string) error {
	var client models.Client

	if err := s.db.Where("identifier = ?", identifier).First(&client).Error; err != nil {
		return err
	}

	if err := s.db.Delete(&client).Error; err != nil {
		return err
	}

	return nil
}

func (s *authService) VerifyClient(identifier, secret string) bool {
	var client models.Client

	if err := s.db.Where("identifier = ?", identifier).First(&client).Error; err != nil {
		return false
	}

	return client.Secret == secret
}

func (s *authService) CreateGroup(group *schemas.GroupCreate) (*schemas.GroupResponse, error) {
	groupModel := schemas.GroupFromCreate(group)

	if err := s.db.Create(groupModel).Error; err != nil {
		return nil, err
	}

	returnGroup := schemas.GroupResponseFromModel(groupModel)

	return returnGroup, nil
}

func (s *authService) GetGroup(identifier string) (*schemas.GroupResponse, error) {
	var group models.Group

	if err := s.db.Where("identifier = ?", identifier).First(&group).Error; err != nil {
		return nil, err
	}

	returnGroup := schemas.GroupResponseFromModel(&group)

	return returnGroup, nil
}

func (s *authService) UpdateGroup(identifier string, group *schemas.GroupUpdate) (*schemas.GroupResponse, error) {
	var existing models.Group

	if err := s.db.Where("identifier = ?", identifier).First(&existing).Error; err != nil {
		return nil, err
	}

	updateData := utils.MakeObjectWithoutNilFields(group)

	if len(updateData) == 0 {
		return schemas.GroupResponseFromModel(&existing), nil
	}

	if err := s.db.Model(&existing).Updates(updateData).Error; err != nil {
		return nil, err
	}

	returnGroup := schemas.GroupResponseFromModel(&existing)
	return returnGroup, nil
}

func (s *authService) DeleteGroup(identifier string) error {
	var group models.Group

	if err := s.db.Where("identifier = ?", identifier).First(&group).Error; err != nil {
		return err
	}

	if err := s.db.Delete(&group).Error; err != nil {
		return err
	}

	return nil
}

func (s *authService) CreateToken(token *schemas.TokenCreate) (*schemas.TokenResponse, error) {
	tokenModel := schemas.TokenFromCreate(token)

	if err := s.db.Create(tokenModel).Error; err != nil {
		return nil, err
	}

	// Recarrega o modelo completo do banco (com todos os campos atualizados)
	if err := s.db.Preload(clause.Associations).First(tokenModel, tokenModel.ID).Error; err != nil {
		return nil, err
	}

	returnToken := schemas.TokenResponseFromModel(tokenModel)
	returnToken.User = nil

	return returnToken, nil
}

func (s *authService) GetToken(identifier string) (*schemas.TokenResponse, error) {
	var token models.Token

	if err := s.db.Preload("User").Where("identifier = ?", identifier).First(&token).Error; err != nil {
		return nil, err
	}

	returnToken := schemas.TokenResponseFromModel(&token)

	return returnToken, nil
}

func (s *authService) GetTokenByRefreshToken(refreshToken string) (*models.Token, error) {
	var token models.Token

	if err := s.db.Where("refresh_token = ?", refreshToken).First(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

func (s *authService) GetTokenByAccessToken(accessToken string) (*schemas.TokenResponse, error) {
	var token models.Token

	if err := s.db.Preload("User").Where("access_token = ?", accessToken).First(&token).Error; err != nil {
		return nil, err
	}

	returnToken := schemas.TokenResponseFromModel(&token)

	return returnToken, nil
}

func (s *authService) RevokeToken(access_token string) error {
	var token models.Token

	if err := s.db.Where("access_token = ?", access_token).First(&token).Error; err != nil {
		return err
	}

	if err := s.db.Delete(&token).Error; err != nil {
		return err
	}

	return nil
}

func (s *authService) Authorize(accessToken string) bool {
	var token models.Token

	if err := s.db.Where("access_token = ?", accessToken).First(&token).Error; err != nil {
		return false
	}

	if token.ExpiresIn < 0 {
		return false
	}

	if token.ExpiresIn < int(time.Now().Unix()) {
		return false
	}

	return true
}
