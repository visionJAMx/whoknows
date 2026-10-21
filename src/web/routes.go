package web

import (
	"github.com/gin-gonic/gin"
)

func Initrouter() *gin.Engine{
	r := gin.Default()


	// test.go
	r.GET("/test", handleTest)

	// search.go
	r.GET("/", handleSearchPage)
	r.GET("/api/search", handleApiSearch)

	// auth.go
	r.POST("/api/login", handleApiLogin)
	
	return r
}