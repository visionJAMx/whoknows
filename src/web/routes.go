package web

import (
	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/service"
)

func Initrouter(authService *service.AuthService) *gin.Engine {
	r := gin.Default()

	// test.go
	r.GET("/test", handleTest)

	// search.go
	r.GET("/", handleSearchPage)
	r.GET("/api/search", handleApiSearch)

	// auth.go
	r.POST("/api/login", handleApiLogin)
	r.POST("/api/register", handleApiRegister(authService))

	return r
}
