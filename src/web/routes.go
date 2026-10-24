package web

import (
	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/service"
)

func Initrouter(authService *service.AuthService) *gin.Engine {
	r := gin.Default()

	r.POST("/api/register", handleApiRegister(authService))

	return r
}
