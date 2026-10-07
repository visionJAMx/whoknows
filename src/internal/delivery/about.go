package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AboutPage viser information om projektet.
// @Summary Vis informationssiden
// @Tags Pages
// @Produce html
// @Success 200 {string} string "Informationssiden som HTML"
// @Router /about [get]
func AboutPage(c *gin.Context) {
	c.HTML(http.StatusOK, "about", gin.H{
		"title": "About",
		"error": "",
	})
}
