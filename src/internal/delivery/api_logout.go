package delivery

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LogoutHandler rydder sessionen og returnerer JSON.
// @Summary Log brugeren ud
// @Description Rydder den aktuelle brugers session.
// @Tags Authentication
// @Produce json
// @Success 200 {object} AuthResponse
// @Failure 500 {object} AuthResponse
// @Router /api/logout [get]
func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)

	// Fjern bruger-ID og andre værdier fra sessionen.
	session.Clear()

	// Gem den ryddede session, så browserens sessionscookie opdateres.
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
