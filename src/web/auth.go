package web

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/service"
)

const sessionUserIDKey = "user_id"

// hasSession reports whether the request belongs to a logged-in user
func hasSession(c *gin.Context) bool {
	return sessions.Default(c).Get(sessionUserIDKey) != nil
}

func handleLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login", gin.H{
			"title":      "Login",
			"error":      "",
			"hasSession": hasSession(c),
		})
	}
}

func handleRegisterPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "register", gin.H{
			"title":      "Register",
			"error":      "",
			"hasSession": hasSession(c),
		})
	}
}

func handleAPILogin(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := service.LoginInput{
			Username: c.PostForm("username"),
			Password: c.PostForm("password"),
		}

		userID, err := auth.Login(c.Request.Context(), input)

		// missing fields → 422
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

		// wrong username or password → 401
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, AuthResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid username or password",
			})
			return
		}

		// anything else → 500
		if err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, AuthResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Could not complete login",
			})
			return
		}

		// only the user ID goes in the signed cookie, never the password or hash
		session := sessions.Default(c)
		session.Clear()
		session.Set(sessionUserIDKey, userID)

		if err := session.Save(); err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, AuthResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Could not save login session",
			})
			return
		}

		c.JSON(http.StatusOK, AuthResponseWithRedirect{
			StatusCode:  http.StatusOK,
			Message:     "You were logged in",
			RedirectURL: "/",
		})
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

		userID, err := auth.Register(c.Request.Context(), input)

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

		// create 

		session := sessions.Default(c)
		session.Clear()
		session.Set(sessionUserIDKey, userID)

		if err := session.Save(); err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, AuthResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Could not automically create session after registering",
			})
			return
		}

		c.JSON(http.StatusOK, AuthResponse{
			StatusCode: http.StatusOK,
			Message:    "You were successfully registered",
		})
	}
}

func handleAPILogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()

		// saving the cleared session updates the browser's cookie
		if err := session.Save(); err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, AuthResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Could not complete logout",
			})
			return
		}

		c.JSON(http.StatusOK, AuthResponse{
			StatusCode: http.StatusOK,
			Message:    "You were logged out",
		})
	}
}
