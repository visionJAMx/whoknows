package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/multitemplate"
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
	router.HTMLRender = createRenderer()
	router.Static("/static", "./static")

	//Routes
	router.GET("/login", delivery.LoginPage)

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

func createRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("login", "templates/layout.html", "templates/login.html")
	r.AddFromFiles("search", "templates/layout.html", "templates/search.html")
	return r
}
