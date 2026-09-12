package delivery

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/visionJAMx/whoknows/src/repository"
)

func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"error": "",
	})
}

func Register(c *gin.Context, db *sql.DB) {
	username := strings.TrimSpace(c.PostForm("username"))
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")

	if username == "" || email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Username, email and password are required",
		})
		return
	}

	if len(username) < 3 || len(username) > 50 {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Username must be between 3 and 50 characters",
		})
		return
	}

	if len(password) < 12 {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Password must be at least 12 characters",
		})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Could not create account",
		})
		return
	}

	_, err = repository.CreateUser(
		c.Request.Context(),
		db,
		username,
		email,
		string(passwordHash),
	)
	if err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Could not create account",
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}
