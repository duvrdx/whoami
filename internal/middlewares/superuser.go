package middlewares

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func SuperuserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Código executado ANTES do handler principal
		fmt.Println("Middleware executando...")

		// Chama o próximo middleware/handler
		if err := next(c); err != nil {
			return err
		}

		// Código executado DEPOIS do handler principal
		fmt.Println("Middleware finalizado!")
		return nil
	}
}
