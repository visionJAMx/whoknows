package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginPage viser loginformularen.
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"error": "",
	})
}
