package delivery

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/internal/repository"
)

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
