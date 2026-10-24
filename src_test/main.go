package main

import (
	"whoknows-server/api"
	"whoknows-server/handlers"
	"whoknows-server/server"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// serve static files
	r.Static("/static", "./static")

	// 1. Initialize business logic layer
	businessService := handlers.NewBusinessService()

	// 2. Initialize transport/route layer with the business service injected
	srv := server.NewServer(businessService)

	// 3. Register routes onto Gin
	api.RegisterHandlers(r, srv)

	// 4. map templates with layout

	server.InitTemplates()

	println("WhoKnows modular server running on http://localhost:8080...")
	r.Run(":8080")
}
