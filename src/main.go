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
	"github.com/visionJAMx/whoknows/src/routes"
)

func main(){




router := routes.Initrouter()
router.Run(":8080");



}
