package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/service"
)

func handleWeatherPage(weather *service.WeatherService) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := weather.WeatherPage(c.Request.Context())
		if err != nil {
			c.String(http.StatusBadGateway, "Couldn't fetch weather forecast")
			return
		}

		c.String(http.StatusOK, result)
	}
}

func handleAPIWeather(weather *service.WeatherService) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := weather.Forecast(c.Request.Context())
		if err != nil {
			c.String(http.StatusBadGateway, "Couldn't fetch weather forecast")
			return
		}

		c.String(http.StatusOK, result)
	}
}
