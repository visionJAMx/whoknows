package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/service"
)

func handleSearchPage(search *service.SearchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := service.SearchInput{
			Query:    c.Query("q"),
			Language: c.Query("language"),
		}

		result, err := search.SearchPage(c.Request.Context(), input)

		if err != nil {
			c.String(http.StatusInternalServerError, "Couldn't search pages")
			return
		}

		c.String(http.StatusOK, result)
	}
}

func handleAPISearch(search *service.SearchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := service.SearchInput{
			Query:    c.Query("q"),
			Language: c.Query("language"),
		}

		result, err := search.Search(c.Request.Context(), input)
		if err != nil {
			c.String(http.StatusInternalServerError, "Couldn't search pages")
			return
		}

		c.String(http.StatusOK, result)
	}
}
