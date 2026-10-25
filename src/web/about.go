package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleAboutPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "about", gin.H{
			"title": "About",
			"error": "",
		})
	}
}
