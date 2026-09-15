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
			"query":    q,
			"language": language,
			"results":  []any{},
			"error":    "",
		}

		if q == "" {
			c.HTML(http.StatusOK, "search.html", data)
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
			c.HTML(http.StatusInternalServerError, "search.html", data)
			return
		}

		data["results"] = pages
		c.HTML(http.StatusOK, "search.html", data)
	}
}
