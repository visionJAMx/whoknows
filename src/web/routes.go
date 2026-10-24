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

	// pages
	r.GET("/", handleSearchPage(searchService))
	r.GET("/login", handleLoginPage())
	r.GET("/register", handleRegisterPage())
	r.GET("/about", handleAboutPage())
	r.GET("/weather", handleWeatherPage(weatherService))

	// api
	r.POST("/api/login", handleAPILogin(authService))
	r.POST("/api/register", handleAPIRegister(authService))
	r.GET("/api/logout", handleAPILogout(authService))
	r.GET("/api/search", handleAPISearch(searchService))
	r.GET("/api/weather", handleAPIWeather(weatherService))

	r.GET("/health", handleHealth())

	return r
}
