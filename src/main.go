package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/internal/delivery"
	"github.com/visionJAMx/whoknows/src/internal/repository"
)

func main() {
	databasePath := os.Getenv("DATABASE_PATH")
	if databasePath == "" {
		databasePath = "../data/whoknows.db"
	}

	db, err := repository.Open(databasePath)
	if err != nil {
		log.Fatalf("could not connect to repository: %v", err)
	}
	defer db.Close()

	if err := repository.Initialize(context.Background(), db); err != nil {
		log.Fatalf("could not initialize repository: %v", err)
	}

	router := gin.Default()

	sessionSecret := os.Getenv("SESSION_SECRET")
	if len(sessionSecret) < 32 {
		log.Fatal("SESSION_SECRET must contain at least 32 characters")
	}

	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   gin.Mode() == gin.ReleaseMode,
		SameSite: http.SameSiteLaxMode,
	})

	router.Use(sessions.Sessions("whoknows_session", store))
	router.LoadHTMLGlob("templates/*.html")

	// Routes
	router.GET("/login", delivery.LoginPage)

	loginHandler := delivery.NewLoginHandler(db)
	router.POST("/api/login", loginHandler.APILogin)

	router.GET("/", delivery.SearchPage)

	router.GET("/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
