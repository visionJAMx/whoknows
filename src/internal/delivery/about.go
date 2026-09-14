package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AboutPage(c *gin.Context) {
	c.HTML(http.StatusOK, "about", gin.H{
		"title": "About",
		"error": "",
	})
}
