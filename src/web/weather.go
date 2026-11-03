package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/service"
)

func handleWeatherPage(weather *service.WeatherService) gin.HandlerFunc {
	return func(c *gin.Context) {
		forecast, err := weather.Forecast(c.Request.Context())
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
}

func handleAPIWeather(weather *service.WeatherService) gin.HandlerFunc {
	return func(c *gin.Context) {
		forecast, err := weather.Forecast(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadGateway, WeatherErrorResponse{Error: "Could not fetch weather"})
			return
		}

		c.JSON(http.StatusOK, WeatherResponse{Data: *forecast})
	}
}
