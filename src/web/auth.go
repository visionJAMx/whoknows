package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/service"
)

func handleLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "GET /login")
	}
}

func handleRegisterPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "GET /register")
	}
}

func handleAPILogin(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := service.LoginInput{
			Username: c.PostForm("username"),
			Password: c.PostForm("password"),
		}

		result, err := auth.Login(c.Request.Context(), input)
		if err != nil {
			c.String(http.StatusInternalServerError, "Couldn't log in")
			return
		}

		// TODO: create session

		c.String(http.StatusOK, result)
	}
}

func handleAPIRegister(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {

		// call service layer and handle potential error
		input := service.RegisterInput{
			Username: "testUser1",
			Email:    "testUser1@test.com",
			Password: "123456",
		}

		result, err := auth.Register(c.Request.Context(), input)
		if err != nil {
			c.String(http.StatusInternalServerError, "Couldn't create account")
			return
		}

		c.String(http.StatusCreated, result)
	}
}

func handleAPILogout(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := auth.Logout(c.Request.Context())
		if err != nil {
			c.String(http.StatusInternalServerError, "Couldn't log out")
			return
		}

		// TODO: clear session

		c.String(http.StatusOK, result)
	}
}
