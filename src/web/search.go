package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/service"
)

const defaultLanguage = "en"

func handleSearchPage(search *service.SearchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := service.SearchInput{
			Query:    c.Query("q"),
			Language: c.Query("language"),
		}

		if input.Language == "" {
			input.Language = defaultLanguage
		}

		pages, err := search.Search(c.Request.Context(), input)

		data := gin.H{
			"title":    "Search",
			"query":    input.Query,
			"language": input.Language,
			"results":  pages,
			"error":    "",
		}

		if err != nil {
			_ = c.Error(err)
			data["error"] = "Could not search pages"
			c.HTML(http.StatusInternalServerError, "search", data)
			return
		}

		c.HTML(http.StatusOK, "search", data)
	}
}

func handleAPISearch(search *service.SearchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// q must be in the URL, but may be empty
		query, ok := c.GetQuery("q")
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, RequestValidationError{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Missing required query parameter: q",
			})
			return
		}

		input := service.SearchInput{
			Query:    query,
			Language: c.Query("language"),
		}

		if input.Language == "" {
			input.Language = defaultLanguage
		}

		pages, err := search.Search(c.Request.Context(), input)
		if err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, SearchErrorResponse{
				Error: "Could not search pages",
			})
			return
		}

		c.JSON(http.StatusOK, SearchResponse{Data: pages})
	}
}
