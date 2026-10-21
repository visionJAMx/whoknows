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
	"github.com/visionJAMx/whoknows/src/web"
	"log"
)

func main() {
	router := web.Initrouter()

	if err := godotenvvault.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	router.Run(":8080")
}
