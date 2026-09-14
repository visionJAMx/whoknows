package delivery

import (
    "database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSearchPageReturnsHTML(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.LoadHTMLGlob("../../templates/*.html")
	router.GET("/", SearchPage)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	if !strings.Contains(response.Body.String(), "<h1>Search</h1>") {
		t.Errorf(
			"expected response to contain search heading, got %s",
			response.Body.String(),
		)
	}

}

func TestSearchPageShowsQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.LoadHTMLGlob("../../templates/*.html")
	router.GET("/", SearchPage)

	request := httptest.NewRequest(
		http.MethodGet,
		"/?q=golang",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if !strings.Contains(response.Body.String(), "golang") {
		t.Errorf(
			"expected response to contain query, got %s",
			response.Body.String(),
		)
	}
}

func TestApiSearch(t *testing.T) {
    gin.SetMode(gin.TestMode)

    var db *sql.DB

    router := gin.New()
    router.GET("/api/search", func(c *gin.Context) {
       ApiSearch(c, db)
    })

    request := httptest.NewRequest(
       http.MethodGet,
       "/api/search?q=golang",
       nil,
    )

    response := httptest.NewRecorder()

    router.ServeHTTP(response, request)

    if response.Code != http.StatusOK && response.Code != http.StatusInternalServerError {
       t.Errorf(
          "expected status 200 or 500, got %d",
          response.Code,
       )
    }
}

