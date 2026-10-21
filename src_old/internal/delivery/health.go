package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse beskriver applikationens statussvar.
type HealthResponse struct {
	Status string `json:"status" binding:"required"`
}

// HealthHandler viser, at webserveren kan besvare en forespørgsel.
// @Summary Kontrollér at applikationen svarer
// @Description Kontrollerer HTTP-serveren, ikke databaseforbindelsen.
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}
