package schemas

import (
	"github.com/duvrdx/whoami/internal/models"
	"github.com/duvrdx/whoami/internal/utils"
)

// User schemas
type UserCreate struct {
	Identifier string  `json:"identifier"`
	Password   string  `json:"password"`
	Metadata   *string `json:"metadata,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

type UserUpdate struct {
	Identifier *string `json:"identifier,omitempty"`
	Password   *string `json:"password,omitempty"`
	Metadata   *string `json:"metadata,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

type UserResponse struct {
	ID         uint    `json:"id"`
	Identifier string  `json:"identifier"`
	Metadata   *string `json:"metadata"`
	IsActive   bool    `json:"is_active"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

func UserResponseFromModel(user *models.User) *UserResponse {
	return &UserResponse{
		ID:         user.ID,
		Identifier: user.Identifier,
		Metadata:   &user.Metadata,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func UserFromCreate(user *UserCreate) *models.User {
	if user == nil {
		return nil
	}

	if user.Password == "" || user.Identifier == "" {
		return nil
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil
	}

	userModel := &models.User{
		Identifier: user.Identifier,
		Password:   hashedPassword,
	}

	if user.Metadata != nil {
		userModel.Metadata = *user.Metadata
	} else {
		userModel.Metadata = "{}"
	}

	if user.IsActive != nil {
		userModel.IsActive = *user.IsActive
	} else {
		userModel.IsActive = false
	}

	return userModel
}

func UserFromUpdate(user *UserUpdate) *models.User {
	if user == nil {
		return nil
	}

	userModel := &models.User{}

	if user.Identifier != nil {
		userModel.Identifier = *user.Identifier
	}

	if user.Password != nil {
		userModel.Password = *user.Password
	}

	if user.Metadata != nil {
		userModel.Metadata = *user.Metadata
	}

	if user.IsActive != nil {
		userModel.IsActive = *user.IsActive
	}

	return userModel
}

// Group schemas
type GroupCreate struct {
	Identifier  string  `json:"identifier"`
	Description *string `json:"description,omitempty"`
	Metadata    *string `json:"metadata,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type GroupUpdate struct {
	Identifier  *string `json:"identifier,omitempty"`
	Description *string `json:"description,omitempty"`
	Metadata    *string `json:"metadata,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type GroupResponse struct {
	ID          uint    `json:"id"`
	Identifier  string  `json:"identifier"`
	Description *string `json:"description"`
	Metadata    *string `json:"metadata"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func GroupResponseFromModel(group *models.Group) *GroupResponse {
	return &GroupResponse{
		ID:          group.ID,
		Identifier:  group.Identifier,
		Description: &group.Description,
		Metadata:    &group.Metadata,
		IsActive:    group.IsActive,
		CreatedAt:   group.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   group.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func GroupFromCreate(group *GroupCreate) *models.Group {
	if group == nil {
		return nil
	}

	groupModel := &models.Group{
		Identifier: group.Identifier,
	}

	if group.Description != nil {
		groupModel.Description = *group.Description
	}

	if group.Metadata != nil {
		groupModel.Metadata = *group.Metadata
	} else {
		groupModel.Metadata = "{}"
	}

	if group.IsActive != nil {
		groupModel.IsActive = *group.IsActive
	} else {
		groupModel.IsActive = false
	}

	return groupModel
}

func GroupFromUpdate(group *GroupUpdate) *models.Group {
	if group == nil {
		return nil
	}

	groupModel := &models.Group{}

	if group.Identifier != nil {
		groupModel.Identifier = *group.Identifier
	}

	if group.Description != nil {
		groupModel.Description = *group.Description
	}

	if group.Metadata != nil {
		groupModel.Metadata = *group.Metadata
	}

	if group.IsActive != nil {
		groupModel.IsActive = *group.IsActive
	}

	return groupModel
}

// Client schemas
type ClientCreate struct {
	Identifier string `json:"identifier"`
	Secret     string `json:"secret"`
	Grant      string `json:"grant"`
	IsActive   *bool  `json:"is_active,omitempty"`
}

type ClientUpdate struct {
	Identifier *string `json:"identifier,omitempty"`
	Secret     *string `json:"secret,omitempty"`
	Grant      *string `json:"grant,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

type ClientResponse struct {
	ID         uint   `json:"id"`
	Identifier string `json:"identifier"`
	Grant      string `json:"grant"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func ClientResponseFromModel(client *models.Client) *ClientResponse {
	return &ClientResponse{
		ID:         client.ID,
		Identifier: client.Identifier,
		Grant:      client.Grant,
		IsActive:   client.IsActive,
		CreatedAt:  client.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  client.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ClientFromCreate(client *ClientCreate) *models.Client {
	if client == nil {
		return nil
	}

	clientModel := &models.Client{
		Identifier: client.Identifier,
		Secret:     client.Secret,
		Grant:      client.Grant,
	}

	if client.IsActive != nil {
		clientModel.IsActive = *client.IsActive
	} else {
		clientModel.IsActive = false
	}

	return clientModel
}

func ClientFromUpdate(client *ClientUpdate) *models.Client {
	if client == nil {
		return nil
	}

	clientModel := &models.Client{}

	if client.Identifier != nil {
		clientModel.Identifier = *client.Identifier
	}

	if client.Secret != nil {
		clientModel.Secret = *client.Secret
	}

	if client.Grant != nil {
		clientModel.Grant = *client.Grant
	}

	if client.IsActive != nil {
		clientModel.IsActive = *client.IsActive
	}

	return clientModel
}

// Token schemas
type TokenPasswordGrant struct {
	Identifier   string `json:"identifier"`
	Password     string `json:"password"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type TokenCreate struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	UserID       uint   `json:"user_id"`
	ClientID     uint   `json:"client_id"`
}

type TokenResponse struct {
	ID           uint   `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func TokenResponseFromModel(token *models.Token) *TokenResponse {
	return &TokenResponse{
		ID:           token.ID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	}
}

func TokenFromCreate(token *TokenCreate) *models.Token {
	return &models.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
		UserID:       token.UserID,
		ClientID:     token.ClientID,
	}
}
