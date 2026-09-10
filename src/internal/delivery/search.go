package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchPage(c *gin.Context) {
	c.HTML(http.StatusOK, "search.html", gin.H{})
}
