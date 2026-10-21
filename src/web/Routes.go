package web

import (
	

	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/handlers"
)



func Initrouter() *gin.Engine{
	r := gin.Default()

	r.GET("/api", handlers.Register )
	
	return r
}
