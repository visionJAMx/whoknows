package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginPage viser loginformularen.
// @Summary Vis loginformularen
// @Tags Pages
// @Produce html
// @Success 200 {string} string "Loginformularen som HTML"
// @Router /login [get]
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{
		"title": "login",
		"error": "",
	})
}
