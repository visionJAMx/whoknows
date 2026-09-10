package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchPage(c *gin.Context) {
	q := c.Query("q")

	c.HTML(http.StatusOK, "search.html", gin.H{
		"error": "",
		"query": q,
	})
}
