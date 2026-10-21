package main

import (
	/*
		"context"
		"fmt"
		"log"
		"net/http"
		"os"


		"github.com/dotenv-org/godotenvvault"
		"github.com/gin-contrib/multitemplate"
		"github.com/gin-contrib/sessions"
		"github.com/gin-contrib/sessions/cookie"
		"github.com/gin-gonic/gin"
	*/
	"github.com/dotenv-org/godotenvvault"
	"github.com/visionJAMx/whoknows/src/service"
	"github.com/visionJAMx/whoknows/src/web"
	"log"
)

func main() {
	if err := godotenvvault.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	authService := service.NewAuthService()

	router := web.Initrouter(authService)

	router.Run(":8080")
}
