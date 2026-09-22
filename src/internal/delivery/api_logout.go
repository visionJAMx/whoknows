package delivery

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.SetCookie("auth_token", "", -1, "/", "localhost", false, true)

	c.Redirect(http.StatusSeeOther, "/login")
}
