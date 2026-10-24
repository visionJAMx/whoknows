package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleAboutPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "GET /about")
	}
}
