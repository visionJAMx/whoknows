package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/dotenv-org/godotenvvault"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/rendering"
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

	sessionSecret := os.Getenv("SESSION_SECRET")
	if len(sessionSecret) < 32 {
		log.Fatal("SESSION_SECRET must contain at least 32 characters")
	}

	sessionStore := cookie.NewStore([]byte(sessionSecret))
	sessionStore.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   gin.Mode() == gin.ReleaseMode,
		SameSite: http.SameSiteLaxMode,
	})

	userRepository := repositories.NewUserRepository(db)
	pageRepository := repositories.NewPageRepository(db)

	authService := service.NewAuthService(userRepository)
	searchService := service.NewSearchService(pageRepository)
	weatherService := service.NewWeatherService()

	router := web.InitRouter(sessionStore, authService, searchService, weatherService)
	router.HTMLRender = rendering.CreateRenderer()

	router.Static("/static", "./static")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
