package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	_ "github.com/vagonaizer/effective-mobile/subscription-service/api/swagger"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/app"
)

// @title Subscription Service API
// @version 1.0
// @description REST API for managing user subscriptions
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @schemes http https

// @tag.name subscriptions
// @tag.description Subscription management operations

// @tag.name health
// @tag.description Health check operations

// @tag.name costs
// @tag.description Cost calculation operations

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

const defaultConfigPath = "configs/config.yaml"

func main() {
	printHello()

	configPath := flag.String("config", defaultConfigPath, "path to configuration file")
	flag.Parse()

	if envConfigPath := os.Getenv("CONFIG_PATH"); envConfigPath != "" {
		*configPath = envConfigPath
	}

	application, err := app.New(*configPath)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("application error: %v", err)
	}
}

func printHello() {
	hello := color.CyanString(`
╔═══════════════════════════════════════╗
║         Subscription Service          ║
║    https://github.com/vagonaizer      ║
╚═══════════════════════════════════════╝
	`)

	fmt.Println(hello)
}
