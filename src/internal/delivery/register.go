package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginPage viser loginformularen.
func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register", gin.H{
		"title": "Register",
		"error": "",
	})
}
