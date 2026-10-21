package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/service"
)

func handleApiLogin(c *gin.Context) {
	// call service layer
	// create session

	c.String(http.StatusOK, "Logged in")
}

func handleApiRegister(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {

		// call service layer and handle potential error
		input := service.RegisterInput{
			Username: "testUser1",
			Email:    "testUser1@test.com",
			Password: "123456",
		}

		if err := auth.Register(c.Request.Context(), input); err != nil {
			c.String(http.StatusInternalServerError, "Couldn't create account")
			return
		}

		c.String(http.StatusCreated, "Account created")
	}
}
