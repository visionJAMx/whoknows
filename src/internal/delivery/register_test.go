package delivery

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterPageReturnsHTML(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.LoadHTMLGlob("../../templates/*.html")
	router.GET("/register", RegisterPage)

	request := httptest.NewRequest(
		http.MethodGet,
		"/register",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	if !strings.Contains(response.Body.String(), "<h1>Register</h1>") {
		t.Errorf(
			"expected response to contain register heading, got %s",
			response.Body.String(),
		)
	}
}
