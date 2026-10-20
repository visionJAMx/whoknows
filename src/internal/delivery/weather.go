package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/internal/weather"
)

func WeatherPage(c *gin.Context) {
	forecast, err := weather.GetForecast(c.Request.Context())
	if err != nil {
		c.HTML(http.StatusBadGateway, "weather", gin.H{
			"title": "Weather",
			"error": "Could not fetch weather forecast",
		})
		return
	}

	c.HTML(http.StatusOK, "weather", gin.H{
		"title":    "Weather",
		"forecast": forecast,
		"error":    "",
	})
}

func WeatherAPI(c *gin.Context) {
	forecast, err := weather.GetForecast(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Could not fetch weather",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": forecast,
	})
}
