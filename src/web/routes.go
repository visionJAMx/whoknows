package web

import (
	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/service"
)

func InitRouter(
	authService *service.AuthService,
	searchService *service.SearchService,
	weatherService *service.WeatherService,
) *gin.Engine {
	r := gin.Default()

	// auth
	r.GET("/login", handleLoginPage())
	r.GET("/register", handleRegisterPage())
	r.POST("/api/login", handleAPILogin(authService))
	r.POST("/api/register", handleAPIRegister(authService))
	r.GET("/api/logout", handleAPILogout(authService))

	// search
	r.GET("/", handleSearchPage(searchService))
	r.GET("/api/search", handleAPISearch(searchService))

	// weather
	r.GET("/weather", handleWeatherPage(weatherService))
	r.GET("/api/weather", handleAPIWeather(weatherService))

	// health
	r.GET("/health", handleHealth())

	return r
}
