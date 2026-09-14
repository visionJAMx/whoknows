package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register", gin.H{
		"title": "Register",
		"error": "",
	})
}
