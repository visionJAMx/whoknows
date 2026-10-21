package web

import "github.com/gin-gonic/gin"

func handleApiLogin(c *gin.Context) {
	// call service layer
	// create session

	c.String(200, "Logged in")
}