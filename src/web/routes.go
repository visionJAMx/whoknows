package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/service"
)

func InitRouter(
	sessionStore sessions.Store,
	authService *service.AuthService,
	searchService *service.SearchService,
	weatherService *service.WeatherService,
) *gin.Engine {
	r := gin.Default()

	// middleware only applies to routes registered after it
	r.Use(sessions.Sessions("whoknows_session", sessionStore))

	// auth
	r.GET("/login", handleLoginPage())
	r.GET("/register", handleRegisterPage())
	r.POST("/api/login", handleAPILogin(authService))
	r.POST("/api/register", handleAPIRegister(authService))
	r.GET("/api/logout", handleAPILogout())

	// search
	r.GET("/", handleSearchPage(searchService))
	r.GET("/api/search", handleAPISearch(searchService))

	// weather
	r.GET("/weather", handleWeatherPage(weatherService))    // TODO
	r.GET("/api/weather", handleAPIWeather(weatherService)) // TODO

	// health
	r.GET("/health", handleHealth())

	return r
}
