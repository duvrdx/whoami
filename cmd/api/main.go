package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/duvrdx/whoami/internal/config"
	"github.com/duvrdx/whoami/internal/models"
	"github.com/duvrdx/whoami/internal/routing"
	"github.com/duvrdx/whoami/internal/schemas"
	"github.com/duvrdx/whoami/internal/utils"
)

func main() {
	config.Init()
	config.Connect()

	config.MigrateDB(models.User{}, models.Client{}, models.Group{}, models.Token{},
		models.RBACRole{}, models.RBACPermission{}, models.RBACResourceType{},
		models.RBACResourceIdentifier{}, models.Config{})

	e := routing.Routing.GetRoutes(routing.Routing{})

	e.HideBanner = true
	e.HidePort = true

	if detectFirstRun() {
		if err := firstRun(); err != nil {
			fmt.Println("Error creating superuser:", err)
			return
		}
		fmt.Println("Superuser created successfully.")
	}

	banner()

	e.Logger.Fatal(e.Start(":7777"))
}

func banner() {

	banner := fmt.Sprintf(""+"| Version: %-10s \n"+"| Started: %-30s", "1.0.0", time.Now().Format(time.RFC1123))
	fmt.Println("\033[1;36m")
	fmt.Print(`
| ▄   ▄  ▗▄▖ ▗▖  ▗▖▗▄▄▄▖
| █ ▄ █ ▐▌ ▐▌▐▛▚▞▜▌  █  
| █▄█▄█ ▐▛▀▜▌▐▌  ▐▌  █  
|       ▐▌ ▐▌▐▌  ▐▌▗▄█▄▖
`)
	fmt.Println("|")
	fmt.Println("| WhoAmI - An IAM API to manage AuthN and AuthZ \n| Made by @duvrdx with Go, Echo, GORM and ☕")
	fmt.Println(banner)
	fmt.Println("\033[0m")
}

func detectFirstRun() bool {
	var countKey, countSuperuser int64

	config.DB.Model(&models.Config{}).Where("key = ?", "first_run").Count(&countKey)

	config.DB.Model(&models.User{}).Where("is_admin = ?", true).Count(&countSuperuser)

	if countKey > 0 && countSuperuser > 0 {
		return false
	}

	return true
}

func firstRun() error {
	var countKey int64

	config.DB.Model(&models.Config{}).Where("key = ?", "first_run").Count(&countKey)

	if countKey > 0 {
		return fmt.Errorf("first run already completed")
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("First run detected. Please create a superuser.")
	fmt.Println("Please enter your identifier:")
	identifier, _ := reader.ReadString('\n')
	identifier = strings.TrimSpace(identifier)

	fmt.Println("Please enter your password:")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	err := createSuperuser(identifier, password)
	if err != nil {
		return err
	}

	first_run := models.Config{
		Key:   "first_run",
		Value: strconv.FormatInt(time.Now().Unix(), 10),
	}

	if err := config.DB.Create(&first_run).Error; err != nil {
		return err
	}

	return nil
}

func createSuperuser(identifier, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if identifier == "" {
		return fmt.Errorf("identifier cannot be empty")
	}
	if len(identifier) < 3 {
		return fmt.Errorf("identifier must be at least 3 characters long")
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if len(password) > 100 {
		return fmt.Errorf("password must be less than 100 characters long")
	}
	if len(identifier) > 100 {
		return fmt.Errorf("identifier must be less than 100 characters long")
	}

	// Check if identifier already exists
	var count int64
	config.DB.Model(&models.User{}).Where("identifier = ?", identifier).Count(&count)
	if count > 0 {
		return fmt.Errorf("identifier already exists")
	}

	superuser := models.User{
		Identifier: identifier,
		Password:   hashedPassword,
		IsAdmin:    true,
		IsActive:   true,
	}

	if err := config.DB.Create(&superuser).Error; err != nil {
		return err
	}

	return nil
}

func createClient(client *schemas.ClientCreate) error {

	if client.Secret == "" {
		return fmt.Errorf("secret cannot be empty")
	}
	if client.Identifier == "" {
		return fmt.Errorf("identifier cannot be empty")
	}
	if len(client.Secret) < 8 {
		return fmt.Errorf("secret must be at least 8 characters long")
	}
	if len(client.Secret) > 100 {
		return fmt.Errorf("secret must be less than 100 characters long")
	}
	if len(client.Identifier) > 100 {
		return fmt.Errorf("identifier must be less than 100 characters long")
	}

	// Check if identifier already exists
	var count int64
	config.DB.Model(&models.Client{}).Where("identifier = ?", client.Identifier).Count(&count)

	clientModel := schemas.ClientFromCreate(client)

	if err := config.DB.Create(clientModel).Error; err != nil {
		return err
	}

	return nil
}
