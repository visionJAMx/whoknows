package delivery

import (
    "database/sql"
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

func TestRegister(t *testing.T) {
    gin.SetMode(gin.TestMode)

    var db *sql.DB

    router := gin.New()

    router.POST("/api/register", func(c *gin.Context) {
       Register(c, db)
    })

    body := strings.NewReader(
       "username=testuser&email=test@example.com&password=secret123",
    )

    request := httptest.NewRequest(
       http.MethodPost,
       "/api/register",
       body,
    )

    request.Header.Set(
       "Content-Type",
       "application/x-www-form-urlencoded",
    )

    response := httptest.NewRecorder()

    router.ServeHTTP(response, request)

    if response.Code != http.StatusSeeOther {
       t.Errorf(
          "expected status %d, got %d",
          http.StatusSeeOther,
          response.Code,
       )
    }

    if location := response.Header().Get("Location"); location != "/login" {
       t.Errorf("expected redirect to /login, got %s", location)
    }
}