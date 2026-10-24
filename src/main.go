package main

import (
	"context"
	"log"
	"os"

	"github.com/dotenv-org/godotenvvault"
	"github.com/visionJAMx/whoknows/src/repositories"
	"github.com/visionJAMx/whoknows/src/service"
	"github.com/visionJAMx/whoknows/src/web"
)

func main() {
	if err := godotenvvault.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	databasePath := os.Getenv("DATABASE_PATH")
	if databasePath == "" {
		databasePath = "../data/whoknows.db"
	}

	db, err := repositories.Open(databasePath)
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}
	defer db.Close()

	if err := repositories.Initialize(context.Background(), db); err != nil {
		log.Fatalf("could not initialize database: %v", err)
	}

	userRepository := repositories.NewUserRepository(db)
	pageRepository := repositories.NewPageRepository(db)

	authService := service.NewAuthService(userRepository)
	searchService := service.NewSearchService(pageRepository)
	weatherService := service.NewWeatherService()

	router := web.InitRouter(authService, searchService, weatherService)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
