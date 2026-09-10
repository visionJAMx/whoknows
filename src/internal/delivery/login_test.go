package delivery

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoginPage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.LoadHTMLFiles("../../templates/login.html")
	router.GET("/login", LoginPage)

	request := httptest.NewRequest(http.MethodGet, "/login", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if !strings.Contains(response.Body.String(), `action="/api/login"`) {
		t.Error("expected login form to submit to /api/login")
	}

	if !strings.Contains(response.Body.String(), `name="username"`) {
		t.Error("expected login form to contain a username field")
	}

	if !strings.Contains(response.Body.String(), `name="password"`) {
		t.Error("expected login form to contain a password field")
	}
}
