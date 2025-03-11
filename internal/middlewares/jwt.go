package middlewares

import (
	"strings"

	"github.com/duvrdx/whoami/internal/config"
	"github.com/duvrdx/whoami/internal/services"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var configJWT = echojwt.Config{
	SigningKey:    config.GetJWTSecret(),
	SigningMethod: "HS256",
	Skipper: func(c echo.Context) bool {
		return strings.HasPrefix(c.Path(), "/o")
	},
	BeforeFunc: func(c echo.Context) {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			c.Logger().Info("Token não encontrado no cabeçalho")
		} else {
			c.Logger().Infof("Token recebido: %s", token)
		}
	},
	ErrorHandler: func(c echo.Context, err error) error {
		c.Logger().Errorf("Erro ao validar JWT: %v", err)
		return echo.ErrUnauthorized
	},
	SuccessHandler: func(c echo.Context) {
		token := c.Request().Header.Get("Authorization")
		authService := services.NewAuthService()

		tokenSplit := strings.Split(token, "Bearer ")
		if len(tokenSplit) != 2 {
			c.Logger().Error("Token inválido")
			c.JSON(401, "Token inválido")
			return
		}

		token = tokenSplit[1]

		tokenInDb, err := authService.GetTokenByAccessToken(token)

		if err != nil {
			c.Logger().Errorf("Erro ao buscar token no banco de dados: %v", err)
			c.JSON(401, "Token inválido")
			return
		}

		if tokenInDb == nil {
			c.Logger().Error("Token não encontrado")
			c.JSON(401, "Token inválido")
			return
		}

		c.Logger().Infof("Token válido")
	},
}

var JWTMiddleware = echojwt.WithConfig(configJWT)
