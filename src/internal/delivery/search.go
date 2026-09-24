package delivery

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/internal/repository"
)

// SearchResponse beskriver søgeresultaterne i OpenAPI.
// Data er en tom liste, når søgningen ikke giver resultater.
type SearchResponse struct {
	Data []map[string]any `json:"data" binding:"required"`
}

// SearchErrorResponse beskriver svaret ved en intern søgefejl.
type SearchErrorResponse struct {
	Error string `json:"error"`
}

// RequestValidationError beskriver kontraktens fejl ved manglende søgeparameter.
type RequestValidationError struct {
	StatusCode int    `json:"statusCode" default:"422"`
	Message    string `json:"message"`
}

// SearchPage viser søgesiden og eventuelle søgeresultater.
// @Summary Vis søgesiden
// @Description Viser søgeformularen og eventuelle søgeresultater.
// @Tags Pages
// @Produce html
// @Param q query string false "Søgetekst"
// @Param language query string false "Sprogfilter" default(en)
// @Success 200 {string} string "Søgesiden som HTML"
// @Failure 500 {string} string "Søgesiden med en fejlbesked"
// @Router / [get]
func SearchPage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Query("q")
		language := c.DefaultQuery("language", "en")

		data := gin.H{
			"query":          q,
			"language":       language,
			"search_results": []any{},
			"error":          "",
		}

		if q == "" {
			c.HTML(http.StatusOK, "search", data)
			return
		}

		pages, err := repository.SearchPages(
			c.Request.Context(),
			db,
			q,
			language,
		)
		if err != nil {
			log.Printf("search pages: %v", err)
			data["error"] = "Could not search pages"
			c.HTML(http.StatusInternalServerError, "search", data)
			return
		}

		data["search_results"] = pages
		c.HTML(http.StatusOK, "search", data)
	}
}

// SearchAPI returnerer søgeresultater som JSON.
// @Summary Søg efter sider
// @Description q skal være med i URL'en, men må være tom. Tom søgning eller ingen resultater giver {"data":[]}.
// @Tags Search
// @Produce json
// @Param q query string true "Søgetekst; tom værdi er tilladt"
// @Param language query string false "Sprogfilter" default(en)
// @Success 200 {object} SearchResponse
// @Failure 422 {object} RequestValidationError "Query-parameteren q mangler"
// @Failure 500 {object} SearchErrorResponse
// @Router /api/search [get]
func SearchAPI(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// q er påkrævet, men en tom værdi (?q=) er tilladt.
		_, exists := c.Request.URL.Query()["q"]
		if !exists {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"statusCode": http.StatusUnprocessableEntity,
				"message":    "Missing required query parameter: q",
			})
			return
		}

		q := c.Query("q")
		language := c.DefaultQuery("language", "en")

		// En tom søgning returnerer en tom liste.
		if q == "" {
			c.JSON(http.StatusOK, gin.H{
				"data": []any{},
			})
			return
		}

		pages, err := repository.SearchPages(
			c.Request.Context(),
			db,
			q,
			language,
		)
		if err != nil {
			log.Printf("api search: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not search pages",
			})
			return
		}

		// Ingen resultater skal give [] i JSON, ikke null.
		if len(pages) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"data": []any{},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": pages,
		})
	}
}
