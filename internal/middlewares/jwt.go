package middlewares

import (
	"strings"

	"github.com/duvrdx/whoami/internal/config"
	"github.com/duvrdx/whoami/internal/services"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func GetJWTMiddleware() echo.MiddlewareFunc {
	var configJWT = echojwt.Config{
		SigningKey:    config.Config.Token.Secret,
		SigningMethod: "HS256",
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Path(), "/o")
		},
		BeforeFunc: func(c echo.Context) {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				c.Logger().Info("Token is not provided")
				c.JSON(401, "Token is not provided")
				return
			} else {
				c.Logger().Infof("Token recebido: %s", token)
			}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			c.Logger().Errorf("Error: %v", err)
			return echo.ErrUnauthorized
		},
		SuccessHandler: func(c echo.Context) {
			token := c.Request().Header.Get("Authorization")
			authService := services.NewAuthService()

			tokenSplit := strings.Split(token, "Bearer ")
			if len(tokenSplit) != 2 {
				c.Logger().Error("Invalid token format")
				c.JSON(401, "Invalid token format")
				return
			}

			token = tokenSplit[1]

			tokenInDb, err := authService.GetTokenByAccessToken(token)

			if err != nil {
				c.Logger().Errorf("Error searching token in DB: %v", err)
				c.JSON(401, "Invalid token")
				return
			}

			if tokenInDb == nil {
				c.Logger().Error("Token not found")
				c.JSON(401, "Invalid token")
				return
			}

			c.Set("user", tokenInDb.User)
		},
	}

	return echojwt.WithConfig(configJWT)
}
