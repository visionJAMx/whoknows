package web

import "github.com/gin-gonic/gin"

func handleTest(c *gin.Context) {
	// call service layer here:
	c.String(200, "OK")
}