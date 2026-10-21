package web

import "github.com/gin-gonic/gin"

func handleSearchPage(c *gin.Context) {
	c.String(200, "OK");
}

func handleApiSearch(c *gin.Context) {
	// call service layer here:
	c.String(200, "OK")
}