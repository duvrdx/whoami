package main

import (
	"github.com/duvrdx/whoami/internal/config"
	"github.com/duvrdx/whoami/internal/models"
	"github.com/duvrdx/whoami/internal/routing"
)

func main() {

	config.InitConfig()

	config.Connect()
	config.MigrateDB(models.User{}, models.Client{}, models.Group{}, models.Token{},
		models.RBACRole{}, models.RBACPermission{}, models.RBACResourceType{},
		models.RBACResourceIdentifier{})

	e := routing.Routing.GetRoutes(routing.Routing{})

	e.Logger.Fatal(e.Start(":8080"))
}
