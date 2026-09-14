package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchPage(c *gin.Context) {
	q := c.Query("q")

	c.HTML(http.StatusOK, "search", gin.H{
		"title": "Search",
		"error": "",
		"query": q,
	})
}
