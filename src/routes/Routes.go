package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


func Initrouter() *gin.Engine{

	r := gin.Default()

	r.GET("/api", func(ctx *gin.Context) {
		ctx.String(http.StatusAccepted, ("hej"))
	})
	

	return r


	

}
