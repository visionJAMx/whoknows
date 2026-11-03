package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/visionJAMx/whoknows/src/models"
)

type PageRepository struct {
	db *sql.DB
}

func NewPageRepository(db *sql.DB) *PageRepository {
	return &PageRepository{db: db}
}

// Search matches query against title and content. An empty language matches all languages.
func (r *PageRepository) Search(ctx context.Context, query, language string) ([]models.Page, error) {
	searchTerm := "%" + query + "%"

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT title, url, language, last_updated, content
		 FROM pages
		 WHERE (? = '' OR language = ?)
		   AND (title LIKE  ? OR content LIKE ?)
		 ORDER BY title`,
		language,
		language,
		searchTerm,
		searchTerm,
	)
	if err != nil {
		return nil, fmt.Errorf("search pages: %w", err)
	}
	defer rows.Close()

	pages := []models.Page{}

	for rows.Next() {
		var page models.Page

		if err := rows.Scan(
			&page.Title,
			&page.URL,
			&page.Language,
			&page.LastUpdated,
			&page.Content,
		); err != nil {
			return nil, fmt.Errorf("read search result: %w", err)
		}

		pages = append(pages, page)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate search results: %w", err)
	}

	return pages, nil
}
