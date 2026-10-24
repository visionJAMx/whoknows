package server

import (
	"net/http"

	"whoknows-server/api"
	"whoknows-server/handlers"

	"github.com/gin-gonic/gin"
)

// Server implements api.ServerInterface and acts as the transport/routing layer
type Server struct {
	Service *handlers.BusinessService
}

func NewServer(service *handlers.BusinessService) *Server {
	return &Server{Service: service}
}

// --- HTML Page Routes ---

func (s *Server) ServeRootPageGet(c *gin.Context) {

	c.Status(http.StatusOK)

	err := appTemplates["search.html"].ExecuteTemplate(c.Writer, "layout.html", nil)
	if err != nil {
		c.Error(err)
	}

}

func (s *Server) ServeWeatherPageWeatherGet(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("<h1>Weather Page</h1>"))
}

func (s *Server) ServeRegisterPageRegisterGet(c *gin.Context) {

	c.Status(http.StatusOK)

	err := appTemplates["register.html"].ExecuteTemplate(c.Writer, "layout.html", nil)
	if err != nil {
		c.Error(err)
	}
}

func (s *Server) ServeLoginPageLoginGet(c *gin.Context) {

	c.Status(http.StatusOK)

	err := appTemplates["login.html"].ExecuteTemplate(c.Writer, "layout.html", nil)
	if err != nil {
		c.Error(err)
	}
}

// --- API Routes ---

func (s *Server) SearchApiSearchGet(c *gin.Context, params api.SearchApiSearchGetParams) {
	res, err := s.Service.Search(params.Q, params.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (s *Server) WeatherApiWeatherGet(c *gin.Context) {
	res, err := s.Service.GetWeather()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (s *Server) RegisterApiRegisterPost(c *gin.Context) {
	var body api.BodyRegisterApiRegisterPost
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, api.HTTPValidationError{})
		return
	}

	code, msg, err := s.Service.Register(body.Username, body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.AuthResponse{
		StatusCode: &code,
		Message:    &msg,
	})
}

func (s *Server) LoginApiLoginPost(c *gin.Context) {
	var body api.BodyLoginApiLoginPost
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, api.HTTPValidationError{})
		return
	}

	code, msg, err := s.Service.Login(body.Username, body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.AuthResponse{
		StatusCode: &code,
		Message:    &msg,
	})
}

func (s *Server) LogoutApiLogoutGet(c *gin.Context) {
	code := 200
	msg := "Logged out successfully"
	c.JSON(http.StatusOK, api.AuthResponse{
		StatusCode: &code,
		Message:    &msg,
	})
}
