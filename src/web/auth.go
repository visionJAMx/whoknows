package web

import (
	"errors"
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
		input := service.RegisterInput{
			Username:        c.PostForm("username"),
			Email:           c.PostForm("email"),
			Password:        c.PostForm("password"),
			PasswordConfirm: c.PostForm("password2"),
		}

		err := auth.Register(c.Request.Context(), input)

		// invalid input → 422
		var validationErr *service.ValidationError
		if errors.As(err, &validationErr) {
			var details []ValidationErrorDetail
			for _, f := range validationErr.Fields {
				details = append(details, ValidationErrorDetail{
					Loc:  []string{"body", f.Field},
					Msg:  f.Message,
					Type: errorTypes[f.Problem],
				})
			}

			c.JSON(http.StatusUnprocessableEntity, ValidationErrorResponse{Detail: details})
			return
		}

		// anything else → 500
		if err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, AuthResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Could not create account",
			})
			return
		}

		c.JSON(http.StatusOK, AuthResponse{
			StatusCode: http.StatusOK,
			Message:    "You were successfully registered",
		})
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
