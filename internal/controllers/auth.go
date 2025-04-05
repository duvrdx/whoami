package controllers

import (
	"time"

	"github.com/duvrdx/whoami/internal/config"
	"github.com/duvrdx/whoami/internal/schemas"
	"github.com/duvrdx/whoami/internal/services"
	"github.com/duvrdx/whoami/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserID   uint `json:"user_id"`
	ClientID uint `json:"client_id"`
	jwt.RegisteredClaims
}

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return AuthController{authService: authService}
}

func (controller AuthController) Register(c echo.Context) error {
	var user schemas.UserCreate

	if err := c.Bind(&user); err != nil {
		return c.JSON(400, err)
	}

	if _, err := controller.authService.CreateUser(&user); err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, "User registered successfully!")
}

func (controller AuthController) UpdateUser(c echo.Context) error {
	var user schemas.UserUpdate
	var identifier = c.Param("identifier")

	if err := c.Bind(&user); err != nil {
		return c.JSON(400, err)
	}

	if user.Password != nil && *user.Password == "" {
		return c.JSON(400, "Password cannot be empty")
	}

	if user.Password != nil && *user.Password != "" {
		hashPassword, err := utils.HashPassword(*user.Password)

		if err != nil {
			return c.JSON(400, err)
		}

		user.Password = &hashPassword
	}

	if _, err := controller.authService.UpdateUser(identifier, &user); err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, "User updated successfully!")
}

func (controller AuthController) GetUser(c echo.Context) error {
	var identifier = c.Param("identifier")

	user, err := controller.authService.GetUser(identifier)

	if err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, user)
}

func (controller AuthController) GetUsers(c echo.Context) error {
	user, err := controller.authService.GetUsers()

	if err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, user)
}

func (controller AuthController) DeleteUser(c echo.Context) error {
	var identifier = c.Param("identifier")

	if err := controller.authService.DeleteUser(identifier); err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(204, "User deleted successfully!")
}

// AuthController methods for OAuth2

func (controller AuthController) CreateClient(c echo.Context) error {
	var client schemas.ClientCreate

	if err := c.Bind(&client); err != nil {
		return c.JSON(400, err)
	}

	if _, err := controller.authService.CreateClient(&client); err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, "Client created successfully!")
}

func (controller AuthController) UpdateClient(c echo.Context) error {
	var client schemas.ClientUpdate
	var identifier = c.Param("identifier")

	if err := c.Bind(&client); err != nil {
		return c.JSON(400, err)
	}

	if _, err := controller.authService.UpdateClient(identifier, &client); err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, "Client updated successfully!")
}

func (controller AuthController) GetClient(c echo.Context) error {
	var identifier = c.Param("identifier")

	client, err := controller.authService.GetClient(identifier)

	if err != nil {
		return c.JSON(404, err)
	}

	return c.JSON(200, client)
}

func (controller AuthController) GetClients(c echo.Context) error {

	clients, err := controller.authService.GetClients()

	if err != nil {
		return c.JSON(404, err)
	}

	return c.JSON(200, clients)
}

func (controller AuthController) DeleteClient(c echo.Context) error {
	var identifier = c.Param("identifier")

	if err := controller.authService.DeleteClient(identifier); err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(204, "Client deleted successfully!")
}

func (controller AuthController) Token(c echo.Context) error {
	var grantType = c.FormValue("grant_type")
	var client_id = c.FormValue("client_id")
	var client_secret = c.FormValue("client_secret")

	client, err := controller.authService.GetClient(client_id)

	if err != nil {
		return c.JSON(404, "Client not found or invalid credentials")
	}

	if !controller.authService.VerifyClient(client_id, client_secret) {
		return c.JSON(404, "Client not found or invalid credentials")
	}

	if client.Grant != grantType {
		return c.JSON(400, "Invalid grant type")
	}

	if grantType == "password" {
		var userIdentifier = c.FormValue("username")
		var userPassword = c.FormValue("password")

		user, err := controller.authService.GetUser(userIdentifier)

		if err != nil {
			return c.JSON(404, "User not found or invalid credentials")
		}

		if !controller.authService.CompareUserPassword(userIdentifier, userPassword) {
			return c.JSON(404, "User not found or invalid credentials")
		}

		// Define o tempo de expiração do token
		expiresIn := time.Now().Unix() + int64(config.Config.Token.Expiration)

		// Cria as claims do JWT
		claims := &Claims{
			UserID:   user.ID,
			ClientID: client.ID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Unix(expiresIn, 0)),
			},
		}

		// Gera o token JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		accessToken, err := token.SignedString(config.Config.Token.Secret)

		if err != nil {
			return c.JSON(500, "Failed to generate access token")
		}

		// Gera o refresh_token (mantido como uma string aleatória)
		refreshToken := utils.GenerateRandomString(16)

		// Cria o token no banco de dados
		tokenData := schemas.TokenCreate{
			AccessToken:  accessToken,  // Agora é um JWT
			RefreshToken: refreshToken, // Continua sendo uma string aleatória
			ExpiresIn:    int(expiresIn),
			UserID:       user.ID,
			ClientID:     client.ID,
		}

		tokenResponse, err := controller.authService.CreateToken(&tokenData)

		if err != nil {
			return c.JSON(400, err)
		}

		// Retorna a resposta com o access_token (JWT) e refresh_token
		return c.JSON(200, tokenResponse)
	}

	return c.JSON(400, "Grant type not implemented")
}

func (controller AuthController) RevokeToken(c echo.Context) error {
	var identifier = c.Param("identifier")

	if err := controller.authService.RevokeToken(identifier); err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(204, "Token revoked successfully!")
}

func (controller AuthController) RefreshToken(c echo.Context) error {
	var refreshToken = c.FormValue("refresh_token")

	token, err := controller.authService.GetTokenByRefreshToken(refreshToken)

	if err != nil {
		return c.JSON(404, "Token not found or invalid")
	}

	if token.ExpiresIn < int(time.Now().Unix()) {
		return c.JSON(404, "Token expired")
	}

	expiresIn := time.Now().Unix() + int64(config.Config.Token.Expiration)

	// Cria as claims do JWT
	claims := &Claims{
		UserID:   token.UserID,
		ClientID: token.ClientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresIn, 0)),
		},
	}

	// Gera o token JWT
	newTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newAccessToken, err := newTokenJWT.SignedString(config.Config.Token.Secret)

	if err != nil {
		return c.JSON(500, "Failed to generate access token")
	}

	newToken := schemas.TokenCreate{
		AccessToken:  newAccessToken,
		RefreshToken: utils.GenerateRandomString(16),
		ExpiresIn:    int(expiresIn),
		UserID:       token.UserID,
		ClientID:     token.ClientID,
	}

	err = controller.authService.RevokeToken(token.AccessToken)

	if err != nil {
		return c.JSON(400, err)
	}

	newTokenResponse, err := controller.authService.CreateToken(&newToken)

	if err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, newTokenResponse)
}

func (controller AuthController) Authorize(c echo.Context) error {
	var accessToken = c.Request().Header.Get("Authorization")

	if accessToken == "" {
		return c.JSON(401, "Unauthorized")
	}

	if !controller.authService.Authorize(accessToken) {
		return c.JSON(401, "Unauthorized")
	}

	return c.JSON(200, "Authorized")
}
